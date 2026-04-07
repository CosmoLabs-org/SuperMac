# Bookmarks

Archived online content - valuable posts, articles, threads, and references worth preserving.

## Purpose

Capture insights from the web before they disappear or become hard to find. Each bookmark preserves:
- Original URL and attribution
- Full content (not just a link)
- Key takeaways for quick reference
- Tags for discoverability

## Organization

```
docs/bookmarks/
├── README.md              # This file
├── x-posts/               # X/Twitter threads
├── blogs/                 # Blog articles
├── tutorials/             # How-to guides
├── tools/                 # Tool discoveries
└── YYYY-MM-DD-slug.md     # Uncategorized (flat)
```

**Organize by**: Platform, topic, or keep flat - whatever works for your project.

Create subfolders as needed: `/bookmark-this <url> --folder x-posts`

## File Format

**Naming**: `YYYY-MM-DD-slug.md`

| Field | Required | Description |
|-------|----------|-------------|
| Title | Yes | Descriptive title |
| URL | Yes | Original source |
| Author | Yes | Creator attribution |
| Platform | Yes | Source platform (x-post, blog, hackernews, etc.) |
| Date Captured | Yes | When you saved it |
| Tags | No | Comma-separated for search |
| Content | Yes | Full text capture |
| Key Takeaways | Yes | 1-3 bullet summary |

## Platform Tags

| Platform | Tag |
|----------|-----|
| X/Twitter | `x-post` |
| LinkedIn | `linkedin` |
| HackerNews | `hackernews` |
| Blog | `blog` |
| Dev.to | `devto` |
| GitHub | `github` |
| YouTube | `youtube` |
| Podcast | `podcast` |

## Usage

```bash
/bookmark-this <url>                    # Fetch and save
/bookmark-this <url> --title "..."      # Custom title
/bookmark-this <url> --tags ai,workflow # With tags
/bookmark-this --manual                 # Paste content manually
```

## Related

- [knowledge-base/](../knowledge-base/) - Permanent technical articles (authored, not archived)
- [prompts/](../prompts/) - Session handoff prompts
