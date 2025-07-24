import finnhub
import datetime as dt

from app.analysis.earnings.types import EarningsCalendarResponse


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
        earnings_calendar = EarningsCalendarResponse(**response)
        print(earnings_calendar)
