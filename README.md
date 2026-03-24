# Courier Service CLI Application

## Overview
This is [Everest Engineering](https://everest.engineering) assessment; a courier service CLI application that calculates the delivery cost, discount and estimated delivery time for packages based on a set of offers (defined from a CSV file) & dispatch vehicle(s) details (number of vehicles, max speed, max load). The application reads package details from standard input, applies relevant calculations and outputs the results to standard output.

## Architecture & Design Decisions

![Sequence Diagram](https://cdn.elyasasmad.com/elyasasmad/ee-sequence-diagram.png)

1. CSV Offer Catalog
CSV file were chosen to define the offer catalog for its simplicity and ease of use. This allows non-technical users to easily add or modify offers using Office applications (Microsoft Excel / WPS Spreadsheet) without needing to change the code. I also considered using [AirTable](https://airtable.com) for a more user-friendly interface, but decided against it to keep the solution offline and self-contained as a CLI application. Maybe in the future, I can add support for fetching offers from an API or a database for more dynamic use cases.

2. Custom Expression Evaluator
A mini expression compiler (lexer -> parser -> AST -> evaluator) was built (with the help of Claude) to handle the offer conditions parsing and evaluation from CSV, instead of hardcoding conditionals or using regex. This enables flexibility so offers can be added or modified with no code changes.

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
`code`: The offer code (e.g., `OFR001`, `OFR002`, etc.)
`discount`: The discount percentage (0-100) (e.g., `10`, `20`, etc.)
`distance`: The distance condition (e.g., `d < 200`, `50 <= d <= 100` etc.)
`weight`: The weight condition (e.g., `w < 100`, `70 <= w <= 200` etc.)

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
I built a mini expression compiler (lexer -> parser -> AST -> evaluator) (with the help of Claude) to handle the offer conditions parsing and evaluation from CSV, instead of hardcoding conditionals or using regex. This enables flexibility so offers can be added or modified with no code changes.
- Assumption: The conditions in the CSV will be simple expressions involving distance (`d`) and weight (`w`) with basic comparison operators (`<`, `<=`, `>`, `>=`) and logical `AND` / `&&`. No complex expressions or `OR` conditions are expected for now. Only 2 variables (`d` and `w`) are expected in the conditions.
- Trade-off: More upfront code than a simple if/else but much more scalable for future offer additions & business rule changes. Also, grammar is easy to extend (e.g.: adding `OR` conditions, more complex expressions, etc.)
- Motivation: I wanted to explore on building a simple compiler in Go and this seemed like a fun opportunity to do so. The package is also made to be reusable and can be extended for my other use cases in the future.
- Next steps: Add more unit tests for the expression evaluator (to handle edge cases) and potentially open source it as a standalone package. Also, add support for other operators (e.g., `OR`, `!=`, etc.) and more complex expressions.

2. Package Dispatch Algorithm
I implemented a brute-force combinator algorithm to find the optimal package combinations while respecting the weight constraints for each vehicle dispatch. However, this algorithm has a time complexity of $O(2^n)$ in the worst case scenario (assuming unlimited weight capacity).
- Assumption: The number of packages per dispatch is expected to be small (e.g., less than 20), which makes the brute-force approach usable without significant performance issues. (At 20 packages, there are 1,048,576 combinations, which is a bit high but still manageable for a CLI application)
- Trade-off: This approach guarantees the most optimal packing per trip (exact solution) at $O(2^n)$, which is acceptable for small $n$. However, for larger $n$ (e.g., 30+ : which would result in 1,073,741,824 combinations), this would become inefficient.
- Next steps: In a real-world scenario with larger inputs, I would implement a more efficient algorithm such as the Dynamic Programming solution to the Knapsack problem which has a time complexity of $O(n \times W)$ where $W$ is the max weight capacity of the vehicle. This solves the problem of the exponential growth in combinations.

3. In-Memory Data Handling
The application processes all input data in-memory, which is suitable for small to medium-sized datasets.
- Assumption: Since this is a CLI application, it's expected that the input data will be of small-medium size that can fit in memory.
- Trade-off: This allows for a simpler code & faster execution for the expected use cases. The cost is that it may not scale well for very large inputs.


## Author
- [Elyas Asmad](https://github.com/ElyasAsmad)