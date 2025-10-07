import yahooFinance from "yahoo-finance2";

export class YahooFinanceService {
  private static _instance: YahooFinanceService;

  private constructor() {
    // Private constructor to prevent direct instantiation
  }

  static get instance(): YahooFinanceService {
    if (!this._instance) {
      this._instance = new YahooFinanceService();
    }
    return this._instance;
  }

  async getQuote(symbol: string) {
    try {
      const quote = await yahooFinance.quote(symbol);
      return quote;
    } catch (error) {
      console.error("Error fetching quote:", error);
      throw error;
    }
  }
}
