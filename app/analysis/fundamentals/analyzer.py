import finnhub
from pydantic import ValidationError
from app.analysis.fundamentals.types import FinnhubMetricsResponse
from app.utils.logger import logger


class FundamentalAnalyzer:
    def __init__(self, api_key: str):
        if not api_key:
            raise ValueError("FINNHUB_API_KEY is not set")
        self.client = finnhub.Client(api_key=api_key)

    def _fetch_metrics(self, symbol: str) -> FinnhubMetricsResponse:
        try:
            response = self.client.company_basic_financials(
                symbol.upper(), 'all')
            metric_data = response.get("metric", {})
            if not metric_data:
                raise ValueError(f"No metrics found for symbol: {symbol}")
            return FinnhubMetricsResponse(**metric_data)
        except ValidationError as ve:
            logger.error(f"Validation error for symbol {symbol}: {ve}")
            return FinnhubMetricsResponse()
        except Exception as e:
            logger.error(f"Error fetching metrics for symbol {symbol}: {e}")
            return FinnhubMetricsResponse()

    def get_profitability_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            # % of revenue remaining after all expenses
            "net_profit_margin_ttm": metrics.net_profit_margin_ttm,
            # % of revenue remaining after COGS
            "gross_margin_ttm": metrics.gross_margin_ttm,
            "operating_margin_ttm": metrics.operating_margin_ttm,  # core business profitability
            # return on equity, how much profit generated with shareholders' equity (how much profit is generated for each dollar of equity)
            "roe_ttm": metrics.roe_ttm,
            # return on assets, how much profit generated with total assets (how much profit is generated for each dollar of assets)
            "roa_ttm": metrics.roa_ttm,
        }

    def get_growth_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            # how rapidly earnings per share have grown over the last 5 years
            "eps_growth_5y": metrics.eps_growth_5y,
            # annualized revenue growth over the last 5 years
            "revenue_growth_5y": metrics.revenue_growth_5y,
            # (free operating cash flow) annualized growth over the last 5 years - cash flow health
            "focf_cagr_5y": metrics.focf_cagr_5y,
            # (compound annual growth rate) of EBITDA over the last 5 years - grows in core profitability
            "ebitda_cagr_5y": metrics.ebitda_cagr_5y,
        }

    def get_valuation_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            # price to book ratio, how much investors are willing to pay for each dollar of book value
            "pb_quarterly": metrics.pb_quarterly,
            # price to sales ratio, how much investors are willing to pay for each dollar of sales
            "ps_ttm": metrics.ps_ttm,
        }

    def get_liquidity_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            # current assets / current liabilities - ability to cover short-term obligations
            "current_ratio_quarterly": metrics.current_ratio_quarterly,
            # current assets - inventory / current liabilities - ability to cover short-term obligations without selling inventory
            "quick_ratio_quarterly": metrics.quick_ratio_quarterly,
        }

    def get_leverage_metrics(self, metrics: FinnhubMetricsResponse) -> dict:
        return {
            # total debt / total equity - how much debt is used to finance the company compared to equity (higher values mean higher financial risk)
            "debt_to_equity_quarterly": metrics.total_debt_to_equity_quarterly,
            # ability to cover interest expenses with earnings before interest and taxes (higher values mean better ability to cover interest expenses)
            "interest_coverage_ttm": metrics.net_interest_coverage_ttm,
        }

    def get_all_metrics(self, symbol: str) -> dict:
        metrics = self._fetch_metrics(symbol)
        return {
            "profitability": self.get_profitability_metrics(metrics),
            "growth": self.get_growth_metrics(metrics),
            "valuation": self.get_valuation_metrics(metrics),
            "liquidity": self.get_liquidity_metrics(metrics),
            "leverage": self.get_leverage_metrics(metrics),
        }

    def format_telegram_msg(self, metrics: dict) -> str:
        msg = ""

        def format_line(label: str, value, suffix: str = "", precision: int = 2) -> str:
            if value is None:
                return ""
            return f"• {label}: {value:.{precision}f}{suffix}\n"

        # Profitability
        profitability = ""
        profitability += format_line("Net Profit Margin",
                                     metrics['profitability']['net_profit_margin_ttm'], "%")
        profitability += format_line("Gross Margin",
                                     metrics['profitability']['gross_margin_ttm'], "%")
        profitability += format_line("Operating Margin",
                                     metrics['profitability']['operating_margin_ttm'], "%")
        profitability += format_line("ROE (Return on Equity)",
                                     metrics['profitability']['roe_ttm'], "%")
        profitability += format_line("ROA (Return on Assets)",
                                     metrics['profitability']['roa_ttm'], "%")
        if profitability:
            msg += "\n<b>📈 Profitability</b>\n" + profitability

        # Growth
        growth = ""
        growth += format_line("EPS Growth (5Y)",
                              metrics['growth']['eps_growth_5y'], "%")
        growth += format_line("Revenue Growth (5Y)",
                              metrics['growth']['revenue_growth_5y'], "%")
        growth += format_line("Free Cash Flow CAGR (5Y)",
                              metrics['growth']['focf_cagr_5y'], "%")
        growth += format_line("EBITDA CAGR (5Y)",
                              metrics['growth'].get('ebitda_cagr_5y'), "%")
        if growth:
            msg += "\n<b>📊 Growth</b>\n" + growth

        # Valuation
        valuation = ""
        valuation += format_line("P/B Ratio",
                                 metrics['valuation']['pb_quarterly'], "x")
        valuation += format_line("P/S Ratio",
                                 metrics['valuation']['ps_ttm'], "x")
        if valuation:
            msg += "\n<b>💰 Valuation</b>\n" + valuation

        # Liquidity
        liquidity = ""
        liquidity += format_line("Current Ratio",
                                 metrics['liquidity']['current_ratio_quarterly'])
        liquidity += format_line("Quick Ratio",
                                 metrics['liquidity']['quick_ratio_quarterly'])
        if liquidity:
            msg += "\n<b>💧 Liquidity</b>\n" + liquidity

        # Leverage
        leverage = ""
        leverage += format_line("Debt/Equity",
                                metrics['leverage']['debt_to_equity_quarterly'])
        leverage += format_line("Interest Coverage",
                                metrics['leverage']['interest_coverage_ttm'], "x")
        if leverage:
            msg += "\n<b>⚖️ Leverage</b>\n" + leverage

        if msg:
            msg = "📊 <b>Fundamental Analysis</b>\n" + msg
        else:
            msg = "No metrics found for the provided symbol."

        return msg.strip()
