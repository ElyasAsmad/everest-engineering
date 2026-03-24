# Courier Service CLI Application

## Overview
This is [Everest Engineering](https://everest.engineering) assessment; a courier service CLI application that calculates the delivery cost, discount and estimated delivery time for packages based on a set of offers (defined from a CSV file) & dispatch vehicle(s) details (number of vehicles, max speed, max load). The application reads package details from standard input, applies relevant calculations and outputs the results to standard output.

## Architecture & Design Decisions

![Sequence Diagram](https://cdn.elyasasmad.com/elyasasmad/ee-sequence-diagram-latest.png)

1. CSV Offer Catalog
<p>CSV file were chosen to define the offer catalog for its simplicity and ease of use. This allows non-technical users to easily add or modify offers using Office applications (Microsoft Excel / WPS Spreadsheet) without needing to change the code. I also considered using [AirTable](https://airtable.com) for a more user-friendly interface, but decided against it to keep the solution offline and self-contained as a CLI application. In the future, I will add support for fetching offers from an API or a database to support more dynamic use cases.</p>

2. Custom Expression Evaluator
<p>A mini expression compiler (lexer -> parser -> AST -> evaluator) was built (with the help of Claude) to handle the offer conditions parsing and evaluation from CSV, instead of hardcoding conditionals or using regex. This enables flexibility so offers can be added or modified with no code changes.</p>

## Setup
1. Ensure you have Go installed on your machine. If you don't have Go installed, you can download it from [https://golang.org/dl/](https://golang.org/dl/).
2. Clone this repository to your local machine:
```bash
git clone https://github.com/ElyasAsmad/everest-engineering.git
```
3. Navigate to the project directory:
```bash
cd everest-engineering
```
4. Install dependencies:
```bash
go mod tidy
```

## Running the Application
This app accepts 1 argument which is the path to the offers CSV file. The CSV file should have the following format:
```csv
code,discount,distance,weight
OFR001,10,d < 200, 70 <= w <= 200
OFR002,20,50 <= d <= 100, w < 100
# add more offers here
```
Explanation of the CSV columns:
| Column | Description | Example |
|--------|-------------|---------|
| `code` | The offer code | `OFR001`, `OFR002` |
| `discount` | The discount percentage (0-100) | `10`, `20` |
| `distance` | The distance condition | `d < 200`, `50 <= d <= 100` |
| `weight` | The weight condition | `w < 100`, `70 <= w <= 200` |

To run the application, use the following command:
```bash
go run ./cmd/app path/to/offers.csv
```

For example, if you have an `offers.csv` file in the current directory, you can run:
```bash
go run ./cmd/app offers.csv
```

This repository also includes a sample `offers.csv` and `input.txt` files in the root directory that is referred from the assessment PDF document.

To simulate the input, you can use the following command:
```bash
cat input.txt | go run ./cmd/app offers.csv
```

---

Alternatively, you can also use the `Makefile` to run the application with the sample input:
```bash
# This will run the application with the sample offers.csv
make run
```

If you want to run the application with the input from `input.txt`:
```bash
make run-mock
```

## Environment Variables
The application uses the following environment variables for configuration:

<!-- table -->
| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| `EE_LOG_LEVEL` | Set the logging level (e.g., DEBUG, INFO, WARN, ERROR, FATAL). | INFO |

## Testing
This project includes unit and integration tests. You can run all the tests easily using the provided `Makefile`.

To run only the unit tests:
```bash
make test
```

To run only the integration tests:
```bash
make test-integration
```

To run all tests (unit and integration) together:
```bash
make test-all
```

To generate the test coverage profile and launch the HTML report:
```bash
make coverage
make coverage-html
```

## Assumptions and Trade-offs
1. Custom Expression Evaluator
<p>I built a mini expression compiler (lexer -> parser -> AST -> evaluator) (with the help of Claude) to handle the offer conditions parsing and evaluation from CSV, instead of hardcoding conditionals or using regex. This enables flexibility so offers can be added or modified with no code changes.</p>

| Aspect | Details |
|--------|---------|
| **Assumption** | Conditions in CSV will be simple expressions with distance (`d`) and weight (`w`), basic comparison operators (`<`, `<=`, `>`, `>=`), and logical `AND`/`&&`. No complex expressions, `OR` conditions, or additional variables expected. |
| **Trade-off** | More upfront code than simple if/else, but significantly more scalable for future offer additions and business rule changes. Grammar is easy to extend (e.g., `OR` conditions, complex expressions). |
| **Motivation** | Explore building a simple compiler in Go. Package designed to be reusable and extensible for future use cases. |
| **Next Steps** | Add unit tests for expression evaluator edge cases. Consider open-sourcing as a standalone package. Add support for additional operators (`OR`, `!=`, etc.) and more complex expressions. |

2. Package Dispatch Algorithm
<p>I implemented a brute-force combinator algorithm to find the optimal package combinations while respecting the weight constraints for each vehicle dispatch. However, this algorithm has a time complexity of $O(2^n)$ in the worst case scenario (assuming unlimited weight capacity).</p>

| Aspect | Details |
|--------|---------|
| **Assumption** | The number of packages per dispatch is expected to be small (e.g., less than 20), which makes the brute-force approach usable without significant performance issues. At 20 packages, there are 1,048,576 combinations, which is acceptable for a CLI application. |
| **Trade-off** | This approach guarantees the most optimal packing per trip (exact solution) at $O(2^n)$, which is acceptable for small $n$. However, for larger $n$ (e.g., 30+: 1,073,741,824 combinations), this becomes inefficient. |
| **Next Steps** | In a real-world scenario with larger inputs, implement a more efficient algorithm such as Dynamic Programming solution to the Knapsack problem with time complexity of $O(n \times W)$ where $W$ is the max weight capacity of the vehicle. |

3. In-Memory Data Handling
<p>The application processes all input data in-memory, which is suitable for small to medium-sized datasets.</p>

| Aspect | Details |
|--------|---------|
| **Assumption** | Since this is a CLI application, it's expected that the input data will be of small-medium size that can fit in memory. |
| **Trade-off** | This allows for a simpler code & faster execution for the expected use cases. The cost is that it may not scale well for very large inputs. |


## Author
- [Elyas Asmad](https://github.com/ElyasAsmad)