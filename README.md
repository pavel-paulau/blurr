N1QL benchmarking tool
----------------------

[![Go Report Card](https://goreportcard.com/badge/github.com/pavel-paulau/nb)](https://goreportcard.com/report/github.com/pavel-paulau/nb)
[![Build Status](https://travis-ci.org/pavel-paulau/nb.svg?branch=master)](https://travis-ci.org/pavel-paulau/nb)
[![Coverage Status](https://coveralls.io/repos/github/pavel-paulau/nb/badge.svg?branch=master)](https://coveralls.io/github/pavel-paulau/nb?branch=master)

Documents
---------

An example of document:

```
{
  "name": "afe067 aa8edc",
  "email": "f42a91@98e71e.com",
  "street": "a043fb2d",
  "city": "b576cd",
  "county": "3b0209",
  "country": "760087",
  "state": "CA",
  "fullState": "District of Columbia",
  "realm": "472999",
  "coins": 98.32,
  "category": 2,
  "achievements": [
    29,
    149,
    149,
    205,
    250
  ],
  "gmtime": [
    1980,
    11,
    4,
    0,
    0,
    0,
    1,
    309,
    0
  ],
  "year": 1988,
  "body": "afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366afe067aa8edcf42a9198e71eb576cd47366"
}
```

Format:

| Field          | Format        | Length  | Range of values | Number of unique values |
|----------------|---------------|---------|-----------------|-------------------------|
| name           | string        | 12      | N/A             | 281,474,976,710,656     |
| email          | string        | 12      | N/A             | 281,474,976,710,656     |
| street         | string        | 8       | N/A             | 4,294,967,296           |
| city           | string        | 6       | N/A             | 16,777,216              |
| county         | string        | 6       | N/A             | 16,777,216              |
| country        | string        | 6       | N/A             | 16,777,216              |
| realm          | string        | 6       | N/A             | 16,777,216              |
| state          | string        | 2       | N/A             | 57                      |
| fullState      | string        | 4 to 24 | N/A             | 57                      |
| coins          | float         | N/A     | [0.1, 655.35]   | 65,535                  |
| category       | integer       | N/A     | [0, 2]          | 3                       |
| achievements   | integer array | 1 to 10 | [0, 511]        | ∞                       |
| gmtime         | integer array | 9       | N/A             | 12                      |
| year           | integer       | N/A     | [1985, 2000]    | 15                      |
| body           | string        | Vary    | N/A             | ∞                       |

Indexes
-------

| Index   | Statement                             |
|---------|---------------------------------------|
|by_email | SELECT * FROM bucket WHERE email='%s' |
