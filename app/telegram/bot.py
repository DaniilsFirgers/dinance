import os
import asyncio
import datetime as dt
from typing import Optional

from aiogram import Bot, Dispatcher, types, F
from aiogram.filters import Command
from aiogram.types import BotCommand
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import StatesGroup, State

from app.analysis.fundamentals.analyzer import FundamentalAnalyzer
from app.analysis.sentiment.analyzer import SentimentAnalyzer
from app.analysis.sentiment.types import InsiderSentimentMessage
from app.telegram.rate_limiter import RateLimiterQueue
from app.utils.logger import logger


class CheckTicker(StatesGroup):
    waiting_for_symbol = State()


class TelegramBot:
    def __init__(self, rate_limiter: RateLimiterQueue, fundamental_analyzer: FundamentalAnalyzer, sentiment_analyzer: SentimentAnalyzer):
        self.token = os.getenv("TELEGRAM_BOT_TOKEN")
        self.bot = Bot(token=self.token)
        self.dp = Dispatcher()
        self.rate_limiter = rate_limiter
        self.fundamental_analyzer = fundamental_analyzer
        self.sentiment_analyzer = sentiment_analyzer

        self.dp.message.register(self.handle_start, Command("start"))
        self.dp.message.register(self.ping, Command("ping"))
        self.dp.message.register(
            self.cmd_check_symbol, Command("check_symbol"))
        self.dp.message.register(
            self.receive_ticker, CheckTicker.waiting_for_symbol)

    async def set_bot_commands(self):
        commands = [
            BotCommand(command="start", description="Start"),
            BotCommand(command="ping", description="Check bot status"),
            BotCommand(command="check_symbol",
                       description="Check stock symbol"),
        ]
        await self.bot.set_my_commands(commands)

    async def cmd_check_symbol(self, message: types.Message, state: FSMContext):
        await message.answer("Please send company ticker you want to check.")
        await state.set_state(CheckTicker.waiting_for_symbol)

    async def receive_ticker(self, message: types.Message, state: FSMContext):
        symbol = message.text.strip()
        await self.answer(message, f"Processing received <b>{symbol}</b> ticker.")

        try:
            metrics = self.fundamental_analyzer.get_all_metrics(symbol)
        except ValueError as ve:
            return await self.rate_limiter.add_request(lambda: message.answer(ve))

        try:
            sentiment = self.sentiment_analyzer.get_insider_sentiment(symbol)
        except ValueError as ve:
            return await self.rate_limiter.add_request(lambda: message.answer(ve))

        await self.answer(message, self.format_analysis_msg(metrics, sentiment))
        await state.clear()

    async def answer(self, message: types.Message, text: str):
        """Send a message with rate limiting."""
        await self.rate_limiter.add_request(lambda: message.answer(text, parse_mode="HTML"))

    def format_analysis_msg(self, metrics: dict, sentiment: Optional[InsiderSentimentMessage]):
        msg = ""

        def format_line(label: str, value: Optional[float], suffix: str = "", precision: int = 2) -> str:
            if value is None:
                return ""
            try:
                return f"• {label}: {float(value):.{precision}f}{suffix}\n"
            except (TypeError, ValueError):
                return ""

        sections = {
            "📈 Profitability": [
                ("Net Profit Margin", metrics["profitability"].get(
                    "net_profit_margin_ttm"), "%"),
                ("Gross Margin", metrics["profitability"].get(
                    "gross_margin_ttm"), "%"),
                ("Operating Margin", metrics["profitability"].get(
                    "operating_margin_ttm"), "%"),
                ("ROE (Return on Equity)",
                 metrics["profitability"].get("roe_ttm"), "%"),
                ("ROA (Return on Assets)",
                 metrics["profitability"].get("roa_ttm"), "%"),
            ],
            "📊 Growth": [
                ("EPS Growth (5Y)", metrics["growth"].get(
                    "eps_growth_5y"), "%"),
                ("Revenue Growth (5Y)", metrics["growth"].get(
                    "revenue_growth_5y"), "%"),
                ("Free Cash Flow CAGR (5Y)",
                 metrics["growth"].get("focf_cagr_5y"), "%"),
                ("EBITDA CAGR (5Y)", metrics["growth"].get(
                    "ebitda_cagr_5y"), "%"),
            ],
            "💰 Valuation": [
                ("P/B Ratio", metrics["valuation"].get("pb_quarterly"), "x"),
                ("P/S Ratio", metrics["valuation"].get("ps_ttm"), "x"),
            ],
            "💧 Liquidity": [
                ("Current Ratio", metrics["liquidity"].get(
                    "current_ratio_quarterly")),
                ("Quick Ratio", metrics["liquidity"].get(
                    "quick_ratio_quarterly")),
            ],
            "⚖️ Leverage": [
                ("Debt/Equity",
                 metrics["leverage"].get("debt_to_equity_quarterly")),
                ("Interest Coverage", metrics["leverage"].get(
                    "interest_coverage_ttm"), "x"),
            ],
        }

        for section, items in sections.items():
            section_text = ""
            for item in items:
                label, value = item[0], item[1]
                suffix = item[2] if len(item) > 2 else ""
                section_text += format_line(label, value, suffix)
            if section_text:
                msg += f"\n<b>{section}</b>\n{section_text}"

        # check if the latest sentiment is at most 3 months old

        if sentiment and (dt.datetime.now() - sentiment.latest.date) <= dt.timedelta(days=90):
            msg += "\n<b>📊 Insider Sentiment</b>\n"
            msg += format_line(f"MSPR ({sentiment.latest.date.strftime('%Y-%m')})",
                               sentiment.latest.mspr, precision=1)
            msg += format_line(f"Avg MSPR ({sentiment.avg.from_date.strftime('%Y-%m')} - {sentiment.avg.to_date.strftime('%Y-%m')})",
                               sentiment.avg.mspr, precision=1)
            msg += f"Trend: <b>{sentiment.latest.trend}</b>\n"

        return msg.strip() if msg else "No metrics found for the provided symbol."

    async def ping(self, message: types.Message):
        """Handles the /ping command."""
        text = "Bots is active! ✅"
        await self.rate_limiter.add_request(lambda: message.answer(text))

    async def handle_start(self, message: types.Message):
        """Handles the /start command."""
        text = "Hello! I am a bot that will help you to perform stock analysis."
        await self.rate_limiter.add_request(lambda: message.answer(text))

    async def start_polling(self):
        """Start an asyncio task to poll the bot."""
        await self.set_bot_commands()
        await self._start_polling()

    async def _start_polling(self):
        """Starts polling the bot asynchronously."""
        try:
            await self.dp.start_polling(self.bot)
        except Exception as e:
            logger.error(f"Error in bot polling: {e}")
            await asyncio.sleep(1)
            await self.start_polling()
