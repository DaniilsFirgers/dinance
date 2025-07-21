import os
import asyncio
from dotenv import load_dotenv

from app.analysis.fundamentals.analyzer import FundamentalAnalyzer
from app.telegram.bot import TelegramBot
from app.telegram.rate_limiter import RateLimiterQueue


async def main():
    load_dotenv()
    fundamental_analyzer = FundamentalAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"))
    print(os.getenv("FINNHUB_API_KEY"))
    print(os.getenv("TELEGRAM_BOT_TOKEN"))
    rate_limiter = RateLimiterQueue(
        rate=30, per=1, buffer=0.2)  # 30 requests per second with a buffer of 0.2 seconds
    telegaram_bot = TelegramBot(rate_limiter, fundamental_analyzer)

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
