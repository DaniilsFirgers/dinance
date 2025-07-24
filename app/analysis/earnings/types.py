
from enum import Enum
from typing import List, Optional
import datetime

from pydantic import BaseModel, ConfigDict, Field, field_validator


class EarningsCallTime(str, Enum):
    BEFORE_MARKET_OPEN = "bmo"
    AFTER_MARKET_CLOSE = "amc"
    DURING_MARKET_HOUR = "dmh"


class EarningsCalendarItem(BaseModel):
    date: Optional[datetime.date] = None
    eps_actual: Optional[float] = Field(None, alias="epsActual")
    eps_estimate: Optional[float] = Field(None, alias="epsEstimate")
    hour: Optional[EarningsCallTime] = None
    quarter: Optional[int] = None
    revenue_actual: Optional[int] = Field(None, alias="revenueActual")
    revenue_estimate: Optional[int] = Field(None, alias="revenueEstimate")
    symbol: Optional[str] = None
    year: Optional[int] = None

    @field_validator("hour", mode="before")
    @classmethod
    def allow_empty_string(cls, v):
        if v == "":
            return None
        return v

    model_config = ConfigDict(populate_by_name=True)


class EarningsCalendarResponse(BaseModel):
    earnings_calendar: List[EarningsCalendarItem] = Field(
        ..., alias="earningsCalendar")

    model_config = ConfigDict(populate_by_name=True)
