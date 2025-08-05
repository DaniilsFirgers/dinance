import asyncio
import aio_pika
import json

from analytics.app.utils.logger import logger


class MessageBus:
    def __init__(self, amqp_url: str, request_queue: str, response_queue: str):
        self.amqp_url = amqp_url
        self.request_queue = request_queue
        self.response_queue = response_queue

    async def connect(self):
        self.connection = await aio_pika.connect_robust(self.amqp_url)
        self.channel = await self.connection.channel()
        self.request_queue = await self.channel.declare_queue(self.request_queue, durable=True)
        self.result_exchange = await self.channel.declare_exchange('', aio_pika.ExchangeType.DIRECT)

    async def process_message(self, message: aio_pika.IncomingMessage):
        async with message.process():
            data = json.loads(message.body.decode())

            processed_data = await self.handle_task(data)

            if processed_data is not None:
                await self.send_result(processed_data)

    async def handle_task(self, data: dict):
        print(f"Received message: {data}")
        return {"status": "done", "input": data}

    async def send_result(self, processed_data: dict):
        body = json.dumps(processed_data).encode()
        await self.channel.default_exchange.publish(
            aio_pika.Message(
                body=body, delivery_mode=aio_pika.DeliveryMode.PERSISTENT),
            routing_key=self.response_queue,
        )

    async def run(self):
        await self.connect()
        logger.info(
            f"Connected to RabbitMQ at {self.amqp_url}. Listening for messages on queue '{self.request_queue.name}'.")
        await self.request_queue.consume(self.process_message)
        await asyncio.Future()
