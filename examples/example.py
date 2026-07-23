"""
Alexa Scraper — Scrapeless LLM Chat Scraper (Python example)

Docs:  https://docs.scrapeless.com/en/llm-chat-scraper/quickstart/introduction/
Token: https://app.scrapeless.com/passport/login?redirect=/quick-start

Run:
    export SCRAPELESS_API_TOKEN="your_api_token"
    pip install requests
    python example.py
"""

import os
import json
import requests

API_URL = "https://api.scrapeless.com/api/v2/scraper/execute"
API_TOKEN = os.environ.get("SCRAPELESS_API_TOKEN", "YOUR_API_TOKEN")

payload = {
    "actor": "scraper.alexa",
    "input": {
        "prompt": "Recommended attractions in New York",
        "country": "US",
    },
    # Optional: receive the result via webhook instead of the sync response.
    # "webhook": {"url": "https://www.your-webhook.com"},
}

headers = {
    "Content-Type": "application/json",
    "x-api-token": API_TOKEN,
}


def main():
    response = requests.post(API_URL, headers=headers, json=payload, timeout=180)
    response.raise_for_status()

    data = response.json()
    result = data.get("task_result", {})

    print("Status:  ", data.get("status"))
    print("Task ID: ", data.get("task_id"))
    print("\nAnswer:\n", result.get("raw_text") or result.get("md_text", ""))

    for ref in result.get("references", []) or []:
        print(f"- {ref.get('title')} -> {ref.get('url')}")

    for product in result.get("products", []) or []:
        print(f"* {product.get('title')} ({product.get('price')}) -> {product.get('url')}")

    # Full response
    print("\nRaw response:\n", json.dumps(data, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
