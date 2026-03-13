"""Tests for ygoprodeck-get-meta skill."""

import sys
from pathlib import Path


# Add scripts directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "ygoprodeck-get-meta" / "scripts"))
from fetch_meta_deck import (
  CATEGORY_MAP,
  API_PAGE_SIZE,
  MetaDeck,
  build_api_url,
  get_default_from_date,
)


class TestMetaDeck:
  """Test MetaDeck class."""

  def test_from_dict_basic(self):
    """Test creating MetaDeck from dictionary."""
    data = {
      "username": "testuser",
      "deck_name": "Test Deck",
      "deck_description": "A test deck",
      "deck_excerpt": "Test Deck excerpt",
      "youtube_link": None,
      "main_deck": '["12345678","87654321"]',
      "extra_deck": '["11111111"]',
      "side_deck": '["22222222"]',
      "edit_date": "1 day ago",
      "deck_views": 100,
      "deckNum": 12345,
      "comments": 5,
      "cover_card": "12345678",
      "deck_price": "50.00",
      "pretty_url": "test-deck-12345",
      "rating": None,
      "ratingq": None,
      "format": "Tournament Meta Decks",
      "common": None,
      "rare": 10,
      "super": 20,
      "ultra": 30,
      "img_url": "https://example.com/image.jpg",
      "submit_date": "2 days ago",
      "tournamentPlayerName": "Player One",
      "tournamentName": "Test Tournament",
      "tournamentPlayerCount": 64,
      "tournamentPlayerCountIsApproximate": 0,
      "tournamentPlacement": "Top 4",
      "userid": 999,
    }

    deck = MetaDeck.from_dict(data)

    assert deck.username == "testuser"
    assert deck.deck_name == "Test Deck"
    assert deck.main_deck == ["12345678", "87654321"]
    assert deck.extra_deck == ["11111111"]
    assert deck.side_deck == ["22222222"]
    assert deck.deckNum == 12345
    assert deck.tournamentPlayerName == "Player One"
    assert deck.tournamentName == "Test Tournament"
    assert deck.tournamentPlacement == "Top 4"

  def test_from_dict_empty_decks(self):
    """Test creating MetaDeck with empty deck lists."""
    data = {
      "username": "testuser",
      "deck_name": "Empty Deck",
      "deck_description": "",
      "deck_excerpt": "",
      "youtube_link": None,
      "main_deck": "[]",
      "extra_deck": "[]",
      "side_deck": "[]",
      "edit_date": "",
      "deck_views": 0,
      "deckNum": 1,
      "comments": 0,
      "cover_card": "",
      "deck_price": "0.00",
      "pretty_url": "",
      "rating": None,
      "ratingq": None,
      "format": "",
      "common": None,
      "rare": 0,
      "super": 0,
      "ultra": 0,
      "img_url": "",
      "submit_date": "",
    }

    deck = MetaDeck.from_dict(data)

    assert deck.main_deck == []
    assert deck.extra_deck == []
    assert deck.side_deck == []

  def test_to_dict(self):
    """Test converting MetaDeck to dictionary."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Test Deck",
      deck_description="A test deck",
      deck_excerpt="Test excerpt",
      youtube_link="https://youtube.com/watch?v=abc",
      main_deck=["12345678"],
      extra_deck=["87654321"],
      side_deck=["11111111"],
      edit_date="1 day ago",
      deck_views=100,
      deckNum=12345,
      comments=5,
      cover_card="12345678",
      deck_price="50.00",
      pretty_url="test-deck-12345",
      rating=4.5,
      ratingq=10,
      format="Tournament Meta Decks",
      common=None,
      rare=10,
      super=20,
      ultra=30,
      img_url="https://example.com/img.jpg",
      submit_date="2 days ago",
    )

    result = deck.to_dict()

    assert result["username"] == "testuser"
    assert result["deck_name"] == "Test Deck"
    assert result["main_deck"] == '["12345678"]'
    assert result["extra_deck"] == '["87654321"]'
    assert result["side_deck"] == '["11111111"]'
    assert result["deckNum"] == 12345

  def test_to_ydk_basic(self):
    """Test converting MetaDeck to YDK format."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Test Deck",
      deck_description="A test deck",
      deck_excerpt="Test excerpt",
      youtube_link=None,
      main_deck=["12345678", "87654321"],
      extra_deck=["11111111"],
      side_deck=["22222222"],
      edit_date="1 day ago",
      deck_views=100,
      deckNum=12345,
      comments=5,
      cover_card="12345678",
      deck_price="50.00",
      pretty_url="test-deck-12345",
      rating=None,
      ratingq=None,
      format="Tournament Meta Decks",
      common=None,
      rare=10,
      super=20,
      ultra=30,
      img_url="https://example.com/img.jpg",
      submit_date="2 days ago",
    )

    ydk = deck.to_ydk("2026-03-13")

    assert "#metadata:deck_name: Test Deck" in ydk
    assert "#metadata:fetch_date: 2026-03-13" in ydk
    assert "#metadata:deck_description: A test deck" in ydk
    assert "#metadata:ygoprodeck_link: https://ygoprodeck.com/deck/test-deck-12345" in ydk
    assert "#main" in ydk
    assert "12345678" in ydk
    assert "87654321" in ydk
    assert "#extra" in ydk
    assert "11111111" in ydk
    assert "!side" in ydk
    assert "22222222" in ydk

  def test_to_ydk_with_tournament(self):
    """Test YDK output includes tournament info."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Tournament Deck",
      deck_description="A tournament deck",
      deck_excerpt="Tournament Deck",
      youtube_link="https://youtube.com/watch?v=abc",
      main_deck=["12345678"],
      extra_deck=[],
      side_deck=[],
      edit_date="1 day ago",
      deck_views=100,
      deckNum=12345,
      comments=5,
      cover_card="12345678",
      deck_price="50.00",
      pretty_url="tournament-deck-12345",
      rating=None,
      ratingq=None,
      format="Tournament Meta Decks",
      common=None,
      rare=10,
      super=20,
      ultra=30,
      img_url="https://example.com/img.jpg",
      submit_date="2 days ago",
      tournamentPlayerName="Player One",
      tournamentName="Test Tournament",
      tournamentPlayerCount=64,
      tournamentPlayerCountIsApproximate=0,
      tournamentPlacement="Top 4",
      userid=999,
    )

    ydk = deck.to_ydk("2026-03-13")

    assert "#metadata:tournament_player_name: Player One" in ydk
    assert "#metadata:tournament_name: Test Tournament" in ydk
    assert "#metadata:tournament_player_count: 64" in ydk

  def test_to_ydk_approximate_count(self):
    """Test YDK output shows approximate indicator."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Tournament Deck",
      deck_description="A tournament deck",
      deck_excerpt="Tournament Deck",
      youtube_link=None,
      main_deck=["12345678"],
      extra_deck=[],
      side_deck=[],
      edit_date="1 day ago",
      deck_views=100,
      deckNum=12345,
      comments=5,
      cover_card="12345678",
      deck_price="50.00",
      pretty_url="tournament-deck-12345",
      rating=None,
      ratingq=None,
      format="Tournament Meta Decks",
      common=None,
      rare=10,
      super=20,
      ultra=30,
      img_url="https://example.com/img.jpg",
      submit_date="2 days ago",
      tournamentPlayerName="Player One",
      tournamentName="Test Tournament",
      tournamentPlayerCount=150,
      tournamentPlayerCountIsApproximate=1,
      tournamentPlacement="Top 4",
      userid=999,
    )

    ydk = deck.to_ydk("2026-03-13")

    assert "#metadata:tournament_player_count: 150 (approximate)" in ydk

  def test_get_safe_filename_with_pretty_url(self):
    """Test safe filename generation with pretty_url."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Test Deck",
      deck_description="",
      deck_excerpt="",
      youtube_link=None,
      main_deck=[],
      extra_deck=[],
      side_deck=[],
      edit_date="",
      deck_views=0,
      deckNum=12345,
      comments=0,
      cover_card="",
      deck_price="0.00",
      pretty_url="my-cool-deck-12345",
      rating=None,
      ratingq=None,
      format="",
      common=None,
      rare=0,
      super=0,
      ultra=0,
      img_url="",
      submit_date="",
    )

    assert deck.get_safe_filename() == "my-cool-deck-12345.ydk"

  def test_get_safe_filename_without_pretty_url(self):
    """Test safe filename generation without pretty_url."""
    deck = MetaDeck(
      username="testuser",
      deck_name="Test Deck: Special!",
      deck_description="",
      deck_excerpt="",
      youtube_link=None,
      main_deck=[],
      extra_deck=[],
      side_deck=[],
      edit_date="",
      deck_views=0,
      deckNum=12345,
      comments=0,
      cover_card="",
      deck_price="0.00",
      pretty_url="",
      rating=None,
      ratingq=None,
      format="",
      common=None,
      rare=0,
      super=0,
      ultra=0,
      img_url="",
      submit_date="",
    )

    filename = deck.get_safe_filename()
    # Regex replaces < > : " / \ | ? * with _
    # So "Test Deck: Special!" becomes "Test Deck_ Special!"
    assert filename == "Test Deck_ Special!_12345.ydk"
    # Verify no invalid characters
    assert "<" not in filename
    assert ">" not in filename
    assert "|" not in filename
    assert "?" not in filename
    assert "*" not in filename


class TestBuildApiUrl:
  """Test build_api_url function."""

  def test_tcg_format(self):
    """Test building URL for TCG format."""
    url = build_api_url("tcg", "2026-01-01", 0)

    assert "_sft_category=Tournament%20Meta%20Decks" in url
    assert "sort=Updated" in url
    assert "from=2026-01-01" in url
    assert "offset=0" in url

  def test_genesys_format(self):
    """Test building URL for Genesys format."""
    url = build_api_url("genesys", "2026-01-01", 0)

    assert "_sft_category=Tournament%20Meta%20Decks%20(Genesys)" in url

  def test_ocg_format(self):
    """Test building URL for OCG format."""
    url = build_api_url("ocg", "2026-01-01", 0)

    assert "_sft_category=Tournament%20Meta%20Decks%20OCG" in url

  def test_default_from_date(self):
    """Test that default from date is used when not provided."""
    url = build_api_url("tcg", None, 0)

    assert "from=" in url
    # Should contain a date (1 month ago from now)
    from_date = get_default_from_date()
    assert from_date in url

  def test_offset_pagination(self):
    """Test offset pagination in URL."""
    url1 = build_api_url("tcg", "2026-01-01", 0)
    url2 = build_api_url("tcg", "2026-01-01", 20)
    url3 = build_api_url("tcg", "2026-01-01", 40)

    assert "offset=0" in url1
    assert "offset=20" in url2
    assert "offset=40" in url3


class TestGetDefaultFromDate:
  """Test get_default_from_date function."""

  def test_returns_correct_format(self):
    """Test that date is in YYYY-MM-DD format."""
    from_date = get_default_from_date()

    # Check format YYYY-MM-DD
    assert len(from_date) == 10
    assert from_date[4] == "-"
    assert from_date[7] == "-"

  def test_is_approximately_one_month_ago(self):
    """Test that date is approximately 30 days ago."""
    from datetime import datetime

    from_date = get_default_from_date()
    parsed = datetime.strptime(from_date, "%Y-%m-%d")
    now = datetime.now()
    diff = now - parsed

    # Should be around 30 days (between 29-31 days)
    assert 29 <= diff.days <= 31


class TestConstants:
  """Test constants."""

  def test_api_page_size(self):
    """Test API page size is 20."""
    assert API_PAGE_SIZE == 20

  def test_category_map(self):
    """Test category map has expected keys."""
    assert "tcg" in CATEGORY_MAP
    assert "genesys" in CATEGORY_MAP
    assert "ocg" in CATEGORY_MAP

  def test_category_map_values(self):
    """Test category map values."""
    assert CATEGORY_MAP["tcg"] == "Tournament%20Meta%20Decks"
    assert CATEGORY_MAP["genesys"] == "Tournament%20Meta%20Decks%20(Genesys)"
    assert CATEGORY_MAP["ocg"] == "Tournament%20Meta%20Decks%20OCG"
