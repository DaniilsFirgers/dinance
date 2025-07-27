from typing import List
import finnhub
import datetime as dt

from app.core.earnings.types import EarningsCalendarItem, EarningsCalendarResponse, EarningsInfo


class EarningsAnalyzer:
    def __init__(self, api_key: str):
        if not api_key:
            raise ValueError("FINNHUB_API_KEY is not set")
        self.client = finnhub.Client(api_key=api_key)
        self.international = False

    def check_calendar(self):
        today = dt.datetime.now().strftime("%Y-%m-%d")
        response = self.client.earnings_calendar(
            _from=today, to=today, symbol="", international=self.international)
        formatted_res = EarningsCalendarResponse(**response)

        valid = self.filter_valid_earnings(formatted_res.earnings_calendar)
        large = self.is_large_enough(valid)

        with_eps_surprise = self.filter_with_eps_surprise(large)

    def format_earnings_info(self, entries: List[EarningsCalendarItem]) -> List[EarningsInfo]:
        return [
            EarningsInfo(symbol=e.symbol,
                         date=e.date,
                         eps_actual=e.eps_actual,
                         eps_estimate=e.eps_estimate,
                         revenue_surprise=e.revenue_actual -
                         e.revenue_estimate if e.revenue_actual and e.revenue_estimate else None,
                         eps_surprise=(e.eps_actual - e.eps_estimate) if e.eps_actual and e.eps_estimate else None)
            for e in entries
        ]

    def is_large_enough(self, entries: List[EarningsCalendarItem], min_revenue=100_000_000) -> List[EarningsCalendarItem]:
        return [
            e for e in entries
            if e.revenue_estimate is not None and e.revenue_estimate >= min_revenue
        ]

    def filter_valid_earnings(self, entries: List[EarningsCalendarItem]) -> List[EarningsCalendarItem]:
        return [
            e for e in entries
            if e.eps_estimate is not None or e.revenue_estimate is not None
        ]

    def filter_with_eps_surprise(self, entries: List[EarningsCalendarItem]) -> List[EarningsCalendarItem]:
        return [
            e for e in entries
            if e.eps_actual is not None and e.eps_estimate is not None and (e.eps_actual > e.eps_estimate or e.eps_actual < e.eps_estimate)
        ]
