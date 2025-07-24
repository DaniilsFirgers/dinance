import os
import asyncio
from dotenv import load_dotenv

from app.analysis.fundamentals.analyzer import FundamentalAnalyzer
from app.analysis.sentiment.analyzer import SentimentAnalyzer
from app.analysis.earnings.analyzer import EarningsAnalyzer
from app.telegram.bot import TelegramBot
from app.telegram.rate_limiter import RateLimiterQueue


async def main():
    load_dotenv()
    fundamental_analyzer = FundamentalAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"))
    sentiment_analyzer = SentimentAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"),
    )

    earnings_calendar = EarningsAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"))

    earnings_calendar.check_calendar()

    rate_limiter = RateLimiterQueue(rate=30, per=1, buffer=0.02)
    telegaram_bot = TelegramBot(
        rate_limiter, fundamental_analyzer, sentiment_analyzer)

    rate_limiter.start()
    # Run bot in a separate task
    polling_task = asyncio.create_task(telegaram_bot.start_polling())

    try:
        await polling_task
    except (KeyboardInterrupt, asyncio.CancelledError):
        print("Shutting down bot...")
    finally:
        await telegaram_bot.bot.session.close()  # 👈 Properly close aiohttp session


if __name__ == "__main__":
    asyncio.run(main())
