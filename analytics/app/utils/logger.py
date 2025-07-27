import logging
import json_log_formatter
import os

log_dir = os.getenv("logs")
os.makedirs("logs", exist_ok=True)


class CustomJsonFormatter(json_log_formatter.JSONFormatter):
    def json_record(self, message, extra, record):
        extra["message"] = message
        extra["level"] = record.levelname
        extra["time"] = self.formatTime(record, self.datefmt)
        return extra


# handle urllib3 logs
urllib_logger = logging.getLogger("urllib3")
urllib_logger.setLevel(logging.WARNING)
urllib_logger.propagate = False  # prevent from propagating to root logger
urllib_logger.handlers.clear()

# handle aiohttp logs
aiohttp_logger = logging.getLogger("aiohttp")
aiohttp_logger.setLevel(logging.WARNING)
aiohttp_logger.propagate = False  # prevent from propagating to root logger
aiohttp_logger.handlers.clear()

# handle sqlachemy logs
sqlalchemy_logger = logging.getLogger("sqlalchemy.engine")
sqlalchemy_logger.setLevel(logging.WARNING)
sqlalchemy_logger.propagate = False  # prevent from propagating to root logger
sqlalchemy_logger.handlers.clear()

logger = logging.getLogger(__name__)

error_log_path = os.path.join("logs", "error.log")
error_file_handler = logging.FileHandler(error_log_path, mode="w")
error_file_handler.setLevel(logging.ERROR)
error_file_handler.setFormatter(CustomJsonFormatter())

# Stream handler (console logging)
stream_handler = logging.StreamHandler()
stream_handler.setFormatter(CustomJsonFormatter())

# Apply base config
logging.basicConfig(
    level=logging.INFO,
    handlers=[
        stream_handler,
        error_file_handler
    ]
)
