# bol CSV Splitter

A small CLI tool to split large **CSV** files into smaller _non-overlapping_ parts.

---

## Installation

**1.** Clone repository.

```bash
git clone https://github.com/fatali-fataliyev/bol-csv-splitter
```

**2.** Build app.

```bash
go build -o bol
```

---

### Example

Example input file is `customers.csv` with 1200 rows.

```bash
bol csv split customers.csv --parts=1,10,rest --out-dir=myFolder
```

You will get these results::

```
customers_part1_1row.csv --> header + 1 row

customers_part2_10rows.csv --> header + next 10 rows

customers_part3_rest_1189_rows.csv --> header + all remaining rows
```

### Rules

- `--parts` flag is required. Example: `1,15,rest`

- `--out-dir` is optional (default: output)

- Existing files will be overwritten.

- Output directory will be created if missing.
