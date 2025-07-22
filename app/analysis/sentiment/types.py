from enum import Enum
from typing import List
from datetime import datetime
from pydantic import BaseModel


class InsiderSentimentTrend(str, Enum):
    STRONG_POSITIVE = "strong_positive"
    POSITIVE = "positive"
    MILD_POSITIVE = "mild_positive"
    NEUTRAL = "neutral"
    MILD_NEGATIVE = "mild_negative"
    NEGATIVE = "negative"
    STRONG_NEGATIVE = "strong_negative"


class LatestSentiment(BaseModel):
    mspr: float
    date: datetime
    trend: InsiderSentimentTrend


class AvgSentiment(BaseModel):
    mspr: float
    from_date: datetime
    to_date: datetime
    trend: InsiderSentimentTrend


class InsiderSentimentMessage(BaseModel):
    avg: AvgSentiment
    latest: LatestSentiment
    delta: float


class InsiderSentimentData(BaseModel):
    symbol: str
    year: int
    month: int
    change: int
    mspr: float  # 100 is most positive, -100 is most negative


class FinnhubInsiderSentimentResponse(BaseModel):
    data: List[InsiderSentimentData]
    symbol: str
