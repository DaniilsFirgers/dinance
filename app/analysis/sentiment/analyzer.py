import datetime as dt
from typing import Optional
import finnhub

from app.analysis.sentiment.types import AvgSentiment, FinnhubInsiderSentimentResponse, InsiderSentimentMessage, InsiderSentimentTrend, LatestSentiment
from app.utils.logger import logger


class SentimentAnalyzer:
    def __init__(self, api_key: str):
        if not api_key:
            raise ValueError("FINNHUB_API_KEY is not set")
        self.client = finnhub.Client(api_key=api_key)

    def get_insider_sentiment(self, symbol: str) -> Optional[float]:
        today = dt.date.today()
        one_year_ago = today - dt.timedelta(days=365)

        start_date = one_year_ago.isoformat()
        end_date = today.isoformat()

        raw_response = self.client.stock_insider_sentiment(
            symbol, start_date, end_date)

        try:
            parsed = FinnhubInsiderSentimentResponse(**raw_response)
            return self._parse_insider_sentiment(parsed)
        except Exception as e:
            logger.error(f"Failed to parse insider sentiment: {e}")
            return None

    def _parse_insider_sentiment(self, response: FinnhubInsiderSentimentResponse, periods: int = 3) -> Optional[InsiderSentimentMessage]:
        if not response.data or len(response.data) == 0:
            return None
        data_sorted = sorted(
            response.data, key=lambda x: dt.datetime(x.year, x.month, 1))

        # Calculate the average sentiment score the the last periods
        selected_periods = data_sorted[-periods:]
        total_mspr = sum(item.mspr for item in selected_periods)
        average_mspr = total_mspr / len(selected_periods)
        last_mspr = selected_periods[-1].mspr

        return InsiderSentimentMessage(
            delta=last_mspr - selected_periods[0].mspr,
            avg=AvgSentiment(
                mspr=average_mspr,
                trend=self._classify_mspr_trend(average_mspr),
                from_date=dt.datetime(
                    selected_periods[0].year, selected_periods[0].month, 1),
                to_date=dt.datetime(
                    selected_periods[-1].year, selected_periods[-1].month, 1)
            ),
            latest=LatestSentiment(
                mspr=last_mspr,
                date=dt.datetime(
                    selected_periods[-1].year, selected_periods[-1].month, 1),
                trend=self._classify_mspr_trend(last_mspr),
            ),
        )

    def _classify_mspr_trend(self, mspr: float) -> InsiderSentimentTrend:
        if mspr > 50:
            return InsiderSentimentTrend.STRONG_POSITIVE
        elif mspr > 30:
            return InsiderSentimentTrend.POSITIVE
        elif mspr > 10:
            return InsiderSentimentTrend.MILD_POSITIVE
        elif mspr > -10:
            return InsiderSentimentTrend.NEUTRAL
        elif mspr > -30:
            return InsiderSentimentTrend.MILD_NEGATIVE
        elif mspr > -50:
            return InsiderSentimentTrend.NEGATIVE
        else:
            return InsiderSentimentTrend.STRONG_NEGATIVE
