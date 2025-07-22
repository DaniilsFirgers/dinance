import finnhub
from typing import Optional, Any

from pydantic import ValidationError
from app.analysis.fundamentals.types import CompanyFinancialsResponse, FinnhubMetricsResponse
from app.utils.logger import logger


class FundamentalAnalyzer:
    def __init__(self, api_key: str):
        if not api_key:
            raise ValueError("FINNHUB_API_KEY is not set")
        self.client = finnhub.Client(api_key=api_key)

    def _fetch_metrics(self, symbol: str) -> FinnhubMetricsResponse:
        response: dict[str, Any] = self.client.company_basic_financials(
            symbol.upper(), 'all')
        try:
            parsed = CompanyFinancialsResponse(**response)
            return parsed.metric
        except ValidationError as ve:
            logger.error(f"Validation error for symbol {symbol}: {ve}")
            return FinnhubMetricsResponse()

    def get_profitability_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            "net_profit_margin_ttm": metrics.net_profit_margin_ttm,
            "gross_margin_ttm": metrics.gross_margin_ttm,
            "operating_margin_ttm": metrics.operating_margin_ttm,
            "roe_ttm": metrics.roe_ttm,
            "roa_ttm": metrics.roa_ttm,
        }

    def get_growth_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            "eps_growth_5y": metrics.eps_growth_5y,
            "revenue_growth_5y": metrics.revenue_growth_5y,
            "focf_cagr_5y": metrics.focf_cagr_5y,
            "ebitda_cagr_5y": metrics.ebitda_cagr_5y,
        }

    def get_valuation_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            "pb_quarterly": metrics.pb_quarterly,
            "ps_ttm": metrics.ps_ttm,
        }

    def get_liquidity_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            "current_ratio_quarterly": metrics.current_ratio_quarterly,
            "quick_ratio_quarterly": metrics.quick_ratio_quarterly,
        }

    def get_leverage_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            "debt_to_equity_quarterly": metrics.total_debt_to_equity_quarterly,
            "interest_coverage_ttm": metrics.net_interest_coverage_ttm,
        }

    def get_all_metrics(self, symbol: str) -> dict:
        try:
            metrics = self._fetch_metrics(symbol)
        except ValueError as e:
            logger.warning(f"Metrics not found for {symbol}: {e}")
            return {}

        return {
            "profitability": self.get_profitability_metrics(metrics),
            "growth": self.get_growth_metrics(metrics),
            "valuation": self.get_valuation_metrics(metrics),
            "liquidity": self.get_liquidity_metrics(metrics),
            "leverage": self.get_leverage_metrics(metrics),
        }

    def format_telegram_msg(self, metrics: dict) -> str:
        msg = ""

        def format_line(label: str, value: Optional[float], suffix: str = "", precision: int = 2) -> str:
            if value is None:
                return ""
            try:
                return f"• {label}: {float(value):.{precision}f}{suffix}\n"
            except (TypeError, ValueError):
                return ""

        sections = {
            "📈 Profitability": [
                ("Net Profit Margin", metrics["profitability"].get(
                    "net_profit_margin_ttm"), "%"),
                ("Gross Margin", metrics["profitability"].get(
                    "gross_margin_ttm"), "%"),
                ("Operating Margin", metrics["profitability"].get(
                    "operating_margin_ttm"), "%"),
                ("ROE (Return on Equity)",
                 metrics["profitability"].get("roe_ttm"), "%"),
                ("ROA (Return on Assets)",
                 metrics["profitability"].get("roa_ttm"), "%"),
            ],
            "📊 Growth": [
                ("EPS Growth (5Y)", metrics["growth"].get(
                    "eps_growth_5y"), "%"),
                ("Revenue Growth (5Y)", metrics["growth"].get(
                    "revenue_growth_5y"), "%"),
                ("Free Cash Flow CAGR (5Y)",
                 metrics["growth"].get("focf_cagr_5y"), "%"),
                ("EBITDA CAGR (5Y)", metrics["growth"].get(
                    "ebitda_cagr_5y"), "%"),
            ],
            "💰 Valuation": [
                ("P/B Ratio", metrics["valuation"].get("pb_quarterly"), "x"),
                ("P/S Ratio", metrics["valuation"].get("ps_ttm"), "x"),
            ],
            "💧 Liquidity": [
                ("Current Ratio", metrics["liquidity"].get(
                    "current_ratio_quarterly")),
                ("Quick Ratio", metrics["liquidity"].get(
                    "quick_ratio_quarterly")),
            ],
            "⚖️ Leverage": [
                ("Debt/Equity",
                 metrics["leverage"].get("debt_to_equity_quarterly")),
                ("Interest Coverage", metrics["leverage"].get(
                    "interest_coverage_ttm"), "x"),
            ],
        }

        for section, items in sections.items():
            section_text = ""
            for item in items:
                label, value = item[0], item[1]
                suffix = item[2] if len(item) > 2 else ""
                section_text += format_line(label, value, suffix)
            if section_text:
                msg += f"\n<b>{section}</b>\n{section_text}"

        return msg.strip() if msg else "No metrics found for the provided symbol."
