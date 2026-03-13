#!/usr/bin/env python3
"""
Fetch meta decks from YGOPRODeck API.

IMPORTANT API BEHAVIOR:
- The YGOPRODeck API always returns exactly 20 items per request.
- To fetch all available decks, we use offset pagination.
- When the response has fewer than 20 items, it means we've reached the end of results.

Usage:
    python fetch_meta_deck.py --format tcg --ydk_dir ./decks     # TCG to YDK files
    python fetch_meta_deck.py --format tcg -o decks.json         # TCG to JSON file
    python fetch_meta_deck.py --format genesys --ydk_dir ./decks # Genesys to YDK
    python fetch_meta_deck.py --format ocg --ydk_dir ./decks     # OCG to YDK
"""

import argparse
import json
import os
import re
import time
import urllib.request
import urllib.parse
from datetime import datetime, timedelta
from typing import Any


API_BASE_URL = "https://ygoprodeck.com/api/decks/getDecks.php"

# API always returns exactly 20 items per request
API_PAGE_SIZE = 20

CATEGORY_MAP = {
  "tcg": "Tournament%20Meta%20Decks",
  "genesys": "Tournament%20Meta%20Decks%20(Genesys)",
  "ocg": "Tournament%20Meta%20Decks%20OCG",
}


class MetaDeck:
  """Represents a meta deck from YGOPRODeck API."""

  def __init__(
    self,
    username: str,
    deck_name: str,
    deck_description: str,
    deck_excerpt: str,
    youtube_link: str | None,
    main_deck: list[str],
    extra_deck: list[str],
    side_deck: list[str],
    edit_date: str,
    deck_views: int,
    deckNum: int,
    comments: int,
    cover_card: str,
    deck_price: str,
    pretty_url: str,
    rating: float | None,
    ratingq: float | None,
    format: str,
    common: int | None,
    rare: int,
    super: int,
    ultra: int,
    img_url: str,
    submit_date: str,
    tournamentPlayerName: str | None = None,
    tournamentName: str | None = None,
    tournamentPlayerCount: int | None = None,
    tournamentPlayerCountIsApproximate: int | None = None,
    tournamentPlacement: str | None = None,
    userid: int | None = None,
  ):
    self.username = username
    self.deck_name = deck_name
    self.deck_description = deck_description
    self.deck_excerpt = deck_excerpt
    self.youtube_link = youtube_link
    self.main_deck = main_deck
    self.extra_deck = extra_deck
    self.side_deck = side_deck
    self.edit_date = edit_date
    self.deck_views = deck_views
    self.deckNum = deckNum
    self.comments = comments
    self.cover_card = cover_card
    self.deck_price = deck_price
    self.pretty_url = pretty_url
    self.rating = rating
    self.ratingq = ratingq
    self.format = format
    self.common = common
    self.rare = rare
    self.super = super
    self.ultra = ultra
    self.img_url = img_url
    self.submit_date = submit_date
    self.tournamentPlayerName = tournamentPlayerName
    self.tournamentName = tournamentName
    self.tournamentPlayerCount = tournamentPlayerCount
    self.tournamentPlayerCountIsApproximate = tournamentPlayerCountIsApproximate
    self.tournamentPlacement = tournamentPlacement
    self.userid = userid

  @classmethod
  def from_dict(cls, data: dict[str, Any]) -> "MetaDeck":
    """Create a MetaDeck from a dictionary."""
    # Parse JSON string fields
    main_deck = json.loads(data.get("main_deck", "[]"))
    extra_deck = json.loads(data.get("extra_deck", "[]"))
    side_deck = json.loads(data.get("side_deck", "[]"))

    return cls(
      username=data.get("username", ""),
      deck_name=data.get("deck_name", ""),
      deck_description=data.get("deck_description", ""),
      deck_excerpt=data.get("deck_excerpt", ""),
      youtube_link=data.get("youtube_link"),
      main_deck=main_deck,
      extra_deck=extra_deck,
      side_deck=side_deck,
      edit_date=data.get("edit_date", ""),
      deck_views=data.get("deck_views", 0),
      deckNum=data.get("deckNum", 0),
      comments=data.get("comments", 0),
      cover_card=data.get("cover_card", ""),
      deck_price=data.get("deck_price", "0.00"),
      pretty_url=data.get("pretty_url", ""),
      rating=data.get("rating"),
      ratingq=data.get("ratingq"),
      format=data.get("format", ""),
      common=data.get("common"),
      rare=data.get("rare", 0),
      super=data.get("super", 0),
      ultra=data.get("ultra", 0),
      img_url=data.get("img_url", ""),
      submit_date=data.get("submit_date", ""),
      tournamentPlayerName=data.get("tournamentPlayerName"),
      tournamentName=data.get("tournamentName"),
      tournamentPlayerCount=data.get("tournamentPlayerCount"),
      tournamentPlayerCountIsApproximate=data.get("tournamentPlayerCountIsApproximate"),
      tournamentPlacement=data.get("tournamentPlacement"),
      userid=data.get("userid"),
    )

  def to_dict(self) -> dict[str, Any]:
    """Convert to dictionary."""
    return {
      "username": self.username,
      "deck_name": self.deck_name,
      "deck_description": self.deck_description,
      "deck_excerpt": self.deck_excerpt,
      "youtube_link": self.youtube_link,
      "main_deck": json.dumps(self.main_deck),
      "extra_deck": json.dumps(self.extra_deck),
      "side_deck": json.dumps(self.side_deck),
      "edit_date": self.edit_date,
      "deck_views": self.deck_views,
      "deckNum": self.deckNum,
      "comments": self.comments,
      "cover_card": self.cover_card,
      "deck_price": self.deck_price,
      "pretty_url": self.pretty_url,
      "rating": self.rating,
      "ratingq": self.ratingq,
      "format": self.format,
      "common": self.common,
      "rare": self.rare,
      "super": self.super,
      "ultra": self.ultra,
      "img_url": self.img_url,
      "submit_date": self.submit_date,
      "tournamentPlayerName": self.tournamentPlayerName,
      "tournamentName": self.tournamentName,
      "tournamentPlayerCount": self.tournamentPlayerCount,
      "tournamentPlayerCountIsApproximate": self.tournamentPlayerCountIsApproximate,
      "tournamentPlacement": self.tournamentPlacement,
      "userid": self.userid,
    }

  def to_ydk(self, fetch_date: str | None = None) -> str:
    """Convert to YDK format string."""
    from datetime import datetime

    # Get current timestamp if not provided
    if fetch_date is None:
      now = datetime.now()
      fetch_date = now.strftime("%Y-%m-%d")
      fetch_date_ts = str(int(now.timestamp()))
    else:
      # Parse the date and convert to timestamp
      dt = datetime.strptime(fetch_date, "%Y-%m-%d")
      fetch_date_ts = str(int(dt.timestamp()))

    lines = []

    # Metadata section
    lines.append(f"#metadata:deck_name: {self.deck_name}")
    lines.append(f"#metadata:fetch_date_ts: {fetch_date_ts}")
    lines.append(f"#metadata:fetch_date: {fetch_date}")
    lines.append(f"#metadata:deck_description: {self.deck_description}")
    lines.append(f"#metadata:deck_excerpt: {self.deck_excerpt}")
    if self.youtube_link:
      lines.append(f"#metadata:youtube_link: {self.youtube_link}")
    lines.append(f"#metadata:ygoprodeck_link: https://ygoprodeck.com/deck/{self.pretty_url}")
    if self.tournamentPlayerName:
      lines.append(f"#metadata:tournament_player_name: {self.tournamentPlayerName}")
    if self.tournamentName:
      lines.append(f"#metadata:tournament_name: {self.tournamentName}")
    if self.tournamentPlayerCount:
      approx = " (approximate)" if self.tournamentPlayerCountIsApproximate else ""
      lines.append(f"#metadata:tournament_player_count: {self.tournamentPlayerCount}{approx}")

    # Main deck
    lines.append("#main")
    for card_id in self.main_deck:
      lines.append(str(card_id))

    # Extra deck
    lines.append("#extra")
    for card_id in self.extra_deck:
      lines.append(str(card_id))

    # Side deck
    lines.append("!side")
    for card_id in self.side_deck:
      lines.append(str(card_id))

    return "\n".join(lines)

  def get_safe_filename(self) -> str:
    """Get a safe filename for the deck."""
    # Use pretty_url if available, otherwise create from deck_name and deckNum
    if self.pretty_url:
      return f"{self.pretty_url}.ydk"
    # Remove invalid filename characters
    safe_name = re.sub(r'[<>:"/\\|?*]', "_", self.deck_name)
    return f"{safe_name}_{self.deckNum}.ydk"


def get_default_from_date() -> str:
  """Get default 'from' date as 1 month ago in YYYY-MM-DD format."""
  one_month_ago = datetime.now() - timedelta(days=30)
  return one_month_ago.strftime("%Y-%m-%d")


def build_api_url(
  category: str,
  from_date: str | None = None,
  offset: int = 0,
) -> str:
  """
  Build the YGOPRODeck API URL for fetching meta decks.

  Note: The API always returns 20 items regardless of limit parameter.
  Use offset to paginate through results.

  Args:
      category: Format category (tcg, genesys, ocg, or custom string)
      from_date: Start date in YYYY-MM-DD format (default: 1 month ago)
      offset: Pagination offset (0, 20, 40, ...)

  Returns:
      Complete API URL string
  """
  # Get category value from map
  category_value = CATEGORY_MAP[category.lower()]

  # Use default from date if not provided
  if from_date is None:
    from_date = get_default_from_date()

  # Note: limit parameter is ignored by the API - it always returns 20 items
  # We include it anyway for documentation purposes
  return f"{API_BASE_URL}?_sft_category={category_value}&sort=Updated&from={from_date}&limit=20&offset={offset}"


def fetch_decks_single_page(url: str) -> list[dict[str, Any]]:
  """
  Fetch a single page of decks from YGOPRODeck API.

  Args:
      url: API URL to fetch from

  Returns:
      List of deck dictionaries (usually 20 items)

  Raises:
      urllib.error.URLError: If request fails
  """
  try:
    # Create request with headers to mimic a browser
    request = urllib.request.Request(
      url,
      headers={
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Accept": "application/json",
        "Accept-Language": "en-US,en;q=0.9",
      },
    )
    with urllib.request.urlopen(request) as response:
      data = response.read().decode("utf-8")
      return json.loads(data)
  except Exception as e:
    print(f"Error fetching decks: {e}")
    return []


def fetch_decks(
  category: str,
  from_date: str | None = None,
) -> list[MetaDeck]:
  """
  Fetch all available decks from YGOPRODeck API with automatic pagination.

  Since the API always returns exactly 20 items per request,
  we fetch all available pages until no more results.

  Args:
      category: Format category (tcg, genesys, ocg, or custom string)
      from_date: Start date in YYYY-MM-DD format (default: 1 month ago)

  Returns:
      List of MetaDeck objects
  """
  all_decks: list[MetaDeck] = []
  offset = 0

  while True:
    url = build_api_url(category, from_date, offset)
    print(f"Fetching page at offset {offset}...")

    deck_dicts = fetch_decks_single_page(url)

    if not deck_dicts:
      # No more results available
      break

    # Convert to MetaDeck objects
    page_decks = [MetaDeck.from_dict(d) for d in deck_dicts]
    all_decks.extend(page_decks)

    # Rate limit: wait 200ms between requests (API supports ~20 req/s)
    time.sleep(0.2)

    # If we got fewer than 20 items, we've reached the end
    if len(page_decks) < API_PAGE_SIZE:
      break

    # Move to next page
    offset += API_PAGE_SIZE

  return all_decks


def main():
  parser = argparse.ArgumentParser(description="Fetch meta decks from YGOPRODeck")
  parser.add_argument(
    "--format",
    "-f",
    default="tcg",
    help="Format: tcg, genesys, ocg, or custom category string (default: tcg)",
  )
  parser.add_argument(
    "--from",
    dest="from_date",
    default=None,
    help="Start date in YYYY-MM-DD format (default: 1 month ago)",
  )
  parser.add_argument(
    "--output",
    "-o",
    default=None,
    help="Store the fetched json in given file path",
  )
  parser.add_argument(
    "--ydk_dir",
    default=None,
    help="Directory to save YDK format deck files",
  )

  args = parser.parse_args()

  # Validate format
  if args.format.lower() not in CATEGORY_MAP:
    parser.error(f"Invalid format: {args.format}. Valid options: {', '.join(CATEGORY_MAP.keys())}")

  # Validate that either --output or --ydk_dir is specified
  if not args.output and not args.ydk_dir:
    parser.error("Either --output or --ydk_dir must be specified")

  print(f"Fetching decks from YGOPRODeck ({args.format} format)...")

  # Fetch all available decks (with automatic pagination)
  decks = fetch_decks(
    category=args.format,
    from_date=args.from_date,
  )

  # If ydk_dir is specified, save decks as YDK files
  if args.ydk_dir:
    from datetime import datetime

    fetch_date = datetime.now().strftime("%Y-%m-%d")
    os.makedirs(args.ydk_dir, exist_ok=True)
    for deck in decks:
      filename = deck.get_safe_filename()
      filepath = os.path.join(args.ydk_dir, filename)
      with open(filepath, "w", encoding="utf-8") as f:
        f.write(deck.to_ydk(fetch_date))
    print(f"Saved {len(decks)} decks to {args.ydk_dir}")
    return

  # Store as JSON file
  decks_dicts = [deck.to_dict() for deck in decks]
  output = json.dumps(decks_dicts, indent=2, ensure_ascii=False)

  with open(args.output, "w", encoding="utf-8") as f:
    f.write(output)
  print(f"Saved {len(decks)} decks to {args.output}")


if __name__ == "__main__":
  main()
