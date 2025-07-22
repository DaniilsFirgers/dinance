import os
import asyncio

from aiogram import Bot, Dispatcher, types, F
from aiogram.filters import Command
from aiogram.types import BotCommand
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import StatesGroup, State

from app.analysis.fundamentals.analyzer import FundamentalAnalyzer
from app.telegram.rate_limiter import RateLimiterQueue
from app.utils.logger import logger


class CheckTicker(StatesGroup):
    waiting_for_symbol = State()


class TelegramBot:
    def __init__(self, rate_limiter: RateLimiterQueue, fundamental_analyzer: FundamentalAnalyzer):
        self.token = os.getenv("TELEGRAM_BOT_TOKEN")
        self.bot = Bot(token=self.token)
        self.dp = Dispatcher()
        self.rate_limiter = rate_limiter
        self.fundamental_analyzer = fundamental_analyzer

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
        try:
            metrics = self.fundamental_analyzer.get_all_metrics(symbol)
        except ValueError as ve:
            await message.answer(ve)
            return
        msg = self.fundamental_analyzer.format_telegram_msg(metrics)
        await self.send_text_msg_with_limiter(msg, message.from_user.id)
        await state.clear()

    async def ping(self, message: types.Message):
        """Handles the /ping command."""
        text = "Bots is active! ✅"
        await self.send_text_msg_with_limiter(text, message.from_user.id)

    async def handle_start(self, message: types.Message):
        """Handles the /start command."""
        text = "Hello! I am a bot that will help you to perform stock analysis."
        await self.send_text_msg_with_limiter(text, message.from_user.id)

    async def send_text_msg_with_limiter(self, message: str, tg_user_id: int):
        """Sends a message to the user."""
        await self.rate_limiter.add_request(lambda: self._send_text_message(message, tg_user_id))

    async def _send_text_message(self, message: str, tg_user_id: int):
        """Sends a message to the user."""
        await self.bot.send_message(
            chat_id=tg_user_id,
            text=message,
            parse_mode="HTML"
        )

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
