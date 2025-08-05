import os
import asyncio
from dotenv import load_dotenv

from app.core.fundamentals.analyzer import FundamentalAnalyzer
from app.core.sentiment.analyzer import SentimentAnalyzer
from app.telegram.bot import TelegramBot
from app.message_bus.rabbitmq import MessageBus
from app.telegram.rate_limiter import RateLimiterQueue


async def main():
    load_dotenv()
    fundamental_analyzer = FundamentalAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"))
    sentiment_analyzer = SentimentAnalyzer(
        api_key=os.getenv("FINNHUB_API_KEY"),
    )

    message_bus = MessageBus(
        amqp_url="amqp://guest:guest@localhost/",
        request_queue="tasks",
        response_queue="results"
    )

    rate_limiter = RateLimiterQueue(rate=30, per=1, buffer=0.02)
    # telegaram_bot = TelegramBot(
    #     rate_limiter, fundamental_analyzer, sentiment_analyzer)

    # rate_limiter.start()
    # Run bot in a separate task
    # polling_task = asyncio.create_task(telegaram_bot.start_polling())

    # try:
    #     await polling_task
    # except (KeyboardInterrupt, asyncio.CancelledError):
    #     print("Shutting down bot...")
    # finally:
    #     await telegaram_bot.bot.session.close()  # 👈 Properly close aiohttp session


if __name__ == "__main__":
    asyncio.run(main())
