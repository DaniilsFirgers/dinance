from typing import Optional
from pydantic import BaseModel, Field


class FinnhubMetricsResponse(BaseModel):
    day_10_avg_trading_volume: Optional[float] = Field(
        None, alias="10DayAverageTradingVolume")
    week_13_price_return_daily: Optional[float] = Field(
        None, alias="13WeekPriceReturnDaily")
    week_26_price_return_daily: Optional[float] = Field(
        None, alias="26WeekPriceReturnDaily")
    month_3_ad_return_std: Optional[float] = Field(
        None, alias="3MonthADReturnStd")
    month_3_avg_trading_volume: Optional[float] = Field(
        None, alias="3MonthAverageTradingVolume")
    week_52_high: Optional[float] = Field(None, alias="52WeekHigh")
    week_52_high_date: Optional[str] = Field(None, alias="52WeekHighDate")
    week_52_low: Optional[float] = Field(None, alias="52WeekLow")
    week_52_low_date: Optional[str] = Field(None, alias="52WeekLowDate")
    week_52_price_return_daily: Optional[float] = Field(
        None, alias="52WeekPriceReturnDaily")
    day_5_price_return_daily: Optional[float] = Field(
        None, alias="5DayPriceReturnDaily")
    asset_turnover_annual: Optional[float] = Field(
        None, alias="assetTurnoverAnnual")
    asset_turnover_ttm: Optional[float] = Field(None, alias="assetTurnoverTTM")
    beta: Optional[float] = Field(None, alias="beta")
    book_value_per_share_annual: Optional[float] = Field(
        None, alias="bookValuePerShareAnnual")
    book_value_per_share_quarterly: Optional[float] = Field(
        None, alias="bookValuePerShareQuarterly")
    capex_cagr_5y: Optional[float] = Field(None, alias="capexCagr5Y")
    cash_flow_per_share_annual: Optional[float] = Field(
        None, alias="cashFlowPerShareAnnual")
    cash_flow_per_share_quarterly: Optional[float] = Field(
        None, alias="cashFlowPerShareQuarterly")
    cash_flow_per_share_ttm: Optional[float] = Field(
        None, alias="cashFlowPerShareTTM")
    cash_per_share_annual: Optional[float] = Field(
        None, alias="cashPerSharePerShareAnnual")
    cash_per_share_quarterly: Optional[float] = Field(
        None, alias="cashPerSharePerShareQuarterly")
    current_dividend_yield_ttm: Optional[float] = Field(
        None, alias="currentDividendYieldTTM")
    current_ratio_annual: Optional[float] = Field(
        None, alias="currentRatioAnnual")
    current_ratio_quarterly: Optional[float] = Field(
        None, alias="currentRatioQuarterly")
    dividend_per_share_ttm: Optional[float] = Field(
        None, alias="dividendPerShareTTM")
    ebitd_per_share_annual: Optional[float] = Field(
        None, alias="ebitdPerShareAnnual")
    ebitd_per_share_ttm: Optional[float] = Field(
        None, alias="ebitdPerShareTTM")
    ebitda_cagr_5y: Optional[float] = Field(None, alias="ebitdaCagr5Y")
    enterprise_value: Optional[float] = Field(None, alias="enterpriseValue")
    eps_annual: Optional[float] = Field(None, alias="epsAnnual")
    eps_basic_excl_extra_items_annual: Optional[float] = Field(
        None, alias="epsBasicExclExtraItemsAnnual")
    eps_basic_excl_extra_items_ttm: Optional[float] = Field(
        None, alias="epsBasicExclExtraItemsTTM")
    eps_excl_extra_items_annual: Optional[float] = Field(
        None, alias="epsExclExtraItemsAnnual")
    eps_excl_extra_items_ttm: Optional[float] = Field(
        None, alias="epsExclExtraItemsTTM")
    eps_growth_3y: Optional[float] = Field(None, alias="epsGrowth3Y")
    eps_growth_5y: Optional[float] = Field(None, alias="epsGrowth5Y")
    eps_incl_extra_items_annual: Optional[float] = Field(
        None, alias="epsInclExtraItemsAnnual")
    eps_incl_extra_items_ttm: Optional[float] = Field(
        None, alias="epsInclExtraItemsTTM")
    eps_normalized_annual: Optional[float] = Field(
        None, alias="epsNormalizedAnnual")
    eps_ttm: Optional[float] = Field(None, alias="epsTTM")
    focf_cagr_5y: Optional[float] = Field(None, alias="focfCagr5Y")
    gross_margin_5y: Optional[float] = Field(None, alias="grossMargin5Y")
    gross_margin_annual: Optional[float] = Field(
        None, alias="grossMarginAnnual")
    gross_margin_ttm: Optional[float] = Field(None, alias="grossMarginTTM")
    inventory_turnover_annual: Optional[float] = Field(
        None, alias="inventoryTurnoverAnnual")
    inventory_turnover_ttm: Optional[float] = Field(
        None, alias="inventoryTurnoverTTM")
    long_term_debt_to_equity_annual: Optional[float] = Field(
        None, alias="longTermDebt/equityAnnual")
    long_term_debt_to_equity_quarterly: Optional[float] = Field(
        None, alias="longTermDebt/equityQuarterly")
    market_capitalization: Optional[float] = Field(
        None, alias="marketCapitalization")
    month_to_date_price_return_daily: Optional[float] = Field(
        None, alias="monthToDatePriceReturnDaily")
    net_income_employee_annual: Optional[float] = Field(
        None, alias="netIncomeEmployeeAnnual")
    net_income_employee_ttm: Optional[float] = Field(
        None, alias="netIncomeEmployeeTTM")
    net_interest_coverage_annual: Optional[float] = Field(
        None, alias="netInterestCoverageAnnual")
    net_interest_coverage_ttm: Optional[float] = Field(
        None, alias="netInterestCoverageTTM")
    net_margin_growth_5y: Optional[float] = Field(
        None, alias="netMarginGrowth5Y")
    net_profit_margin_5y: Optional[float] = Field(
        None, alias="netProfitMargin5Y")
    net_profit_margin_annual: Optional[float] = Field(
        None, alias="netProfitMarginAnnual")
    net_profit_margin_ttm: Optional[float] = Field(
        None, alias="netProfitMarginTTM")
    operating_margin_5y: Optional[float] = Field(
        None, alias="operatingMargin5Y")
    operating_margin_annual: Optional[float] = Field(
        None, alias="operatingMarginAnnual")
    operating_margin_ttm: Optional[float] = Field(
        None, alias="operatingMarginTTM")
    pb: Optional[float] = Field(None, alias="pb")
    pb_annual: Optional[float] = Field(None, alias="pbAnnual")
    pb_quarterly: Optional[float] = Field(None, alias="pbQuarterly")
    ps_annual: Optional[float] = Field(None, alias="psAnnual")
    ps_ttm: Optional[float] = Field(None, alias="psTTM")
    ptbv_annual: Optional[float] = Field(None, alias="ptbvAnnual")
    ptbv_quarterly: Optional[float] = Field(None, alias="ptbvQuarterly")
    quick_ratio_annual: Optional[float] = Field(None, alias="quickRatioAnnual")
    quick_ratio_quarterly: Optional[float] = Field(
        None, alias="quickRatioQuarterly")
    receivables_turnover_annual: Optional[float] = Field(
        None, alias="receivablesTurnoverAnnual")
    receivables_turnover_ttm: Optional[float] = Field(
        None, alias="receivablesTurnoverTTM")
    revenue_employee_annual: Optional[float] = Field(
        None, alias="revenueEmployeeAnnual")
    revenue_employee_ttm: Optional[float] = Field(
        None, alias="revenueEmployeeTTM")
    revenue_growth_3y: Optional[float] = Field(None, alias="revenueGrowth3Y")
    revenue_growth_5y: Optional[float] = Field(None, alias="revenueGrowth5Y")
    revenue_growth_quarterly_yoy: Optional[float] = Field(
        None, alias="revenueGrowthQuarterlyYoy")
    revenue_growth_ttm_yoy: Optional[float] = Field(
        None, alias="revenueGrowthTTMYoy")
    revenue_per_share_annual: Optional[float] = Field(
        None, alias="revenuePerShareAnnual")
    revenue_per_share_ttm: Optional[float] = Field(
        None, alias="revenuePerShareTTM")
    revenue_share_growth_5y: Optional[float] = Field(
        None, alias="revenueShareGrowth5Y")
    roa_5y: Optional[float] = Field(None, alias="roa5Y")
    roa_rfy: Optional[float] = Field(None, alias="roaRfy")
    roa_ttm: Optional[float] = Field(None, alias="roaTTM")
    roe_5y: Optional[float] = Field(None, alias="roe5Y")
    roe_rfy: Optional[float] = Field(None, alias="roeRfy")
    roe_ttm: Optional[float] = Field(None, alias="roeTTM")
    roi_5y: Optional[float] = Field(None, alias="roi5Y")
    roi_annual: Optional[float] = Field(None, alias="roiAnnual")
    roi_ttm: Optional[float] = Field(None, alias="roiTTM")
    tangible_book_value_per_share_annual: Optional[float] = Field(
        None, alias="tangibleBookValuePerShareAnnual")
    tangible_book_value_per_share_quarterly: Optional[float] = Field(
        None, alias="tangibleBookValuePerShareQuarterly")
    total_debt_to_equity_annual: Optional[float] = Field(
        None, alias="totalDebt/totalEquityAnnual")
    total_debt_to_equity_quarterly: Optional[float] = Field(
        None, alias="totalDebt/totalEquityQuarterly")
    year_to_date_price_return_daily: Optional[float] = Field(
        None, alias="yearToDatePriceReturnDaily")

    class Config:
        populate_by_name = True  # Allows using the field name instead of the alias


class CompanyFinancialsResponse(BaseModel):
    metric: FinnhubMetricsResponse
