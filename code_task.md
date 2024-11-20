# Code Task: Data Engineering API Endpoint

This document outlines the details of a REST API endpoint designed to process hierarchical data submitted in CSV format.

## Getting Started with the CSV to JSON Hierarchical Data Converter

This guide will help you get the CSV to JSON converter up and running and perform various tests.

**Before you begin:**
* Make sure you have Go version 1.23 or higher installed on your system.

**Running the Rest API server:**
* Open your terminal and navigate to the project directory.
* Execute the following command:

    ```bash
    go run main.go
    ```
    This will start the http server, which will listen for HTTP requests on port 8080.

**Testing the converter:**

The code includes unit and integration tests to ensure the converter functions correctly.

* **Unit Testing**

    ```bash
    go test -count=1 ./tree/...
    ```
    This command runs unit tests specifically for the core conversion logic. The -count=1 flag disables caching and ensures each test runs independently.

* **Integration Tests:**

    ```bash
    go test -count=1 ./...
    ```
    This command runs all tests, including integration tests that validate the entire conversion process from HTTP request to JSON response. Make sure the REST API is running (as described in the ["Running the Rest API server"](#running-the-rest-api-server) section) before executing this command.

* **Benchmarking**

    The implementation conducted benchmark tests to evaluate the performance of different approaches, including concurrent processing with goroutines. To see these results and understand why the current implementation was chosen, run:
    ```bash
    go test -bench=. ./tree
    ```

## Functional Requirements

The API endpoint handles requests as follows:

**Request:**

* **Method:** `POST`
* **Content-Type:** `text/csv`
* **Body:** A CSV file representing a hierarchy of items in a flat structure. 
    * The file must include a header row with column names.
    * Data rows should be delimited by newline characters (`\n`).
    * Columns should be delimited by commas (`,`).
    * The following columns are required:
        * `item_id`: Unique identifier for each item.
        * `level_1`: Top-level category.
        * `level_2` (optional): Second-level category.
        * `level_3` (optional): Third-level category.
    * Level columns can be empty for leaf nodes.
    * Column order is flexible.

**Example CSV Data:**
- File 1 with item_id as the last column:
```csv
    level_1,level_2,level_3,item_id
    1,12,103,12507622
```
- File 2 with item_id as first column:
```csv
    level_1,level_2,level_3,item_id
    1,12,103,12507622
```
- File 3 with leaf note before lowest level:
```csv
    level_1,level_2,item_id
    category_1,,item_1
    category_2,category_3,item_2
```

**Response:**

* **Content-Type:** `application/json`
* **Status Codes:**
    * **200 OK**: Successful processing. The response body will contain a JSON object representing the hierarchy, like this:
    ```json
    {
        "children": {
            "A": {
                "children": {
                    "1": {
                        "item": true
                    }
                }
            }
        }
    }
    ```
    * **400 Bad Request**:  Indicates an invalid request due to:
        * **Invalid header**:
            * Empty header row.
            * Missing `level_1` or `item_id` columns.
            * Invalid column names (must be one of `item_id`, `level_1`, `level_2`, or `level_3`).
        * **Invalid content:**
            * Mismatch between the number of columns in the header and data rows.
            * "Level skipping" (e.g., a row with a value in `level_2` but no value in `level_1`). Like this:
            ```csv
                level_1,level_2,item_id
                ,category_2,item_1            
            ```
    * **405 Method Not Allowed:** Returned if a method other than POST is used.
    * **415 Unsupported Media Type:** Returned if the Content-Type is not text/csv.

## Non-Functional Requirements

The API endpoint is designed with the following non-functional requirements in mind:

* **Availability:**
    * Ensure high availability through infrastructure redundancy.
* **Scalability:**
    * Handle increasing traffic loads.
    * Accommodate growing data volumes.
* **Performance:**
    * Minimize response times.
    * Maximize throughput.
    * Maintain low latency.
* **Maintainability:**
    * Employ clear and organized code structure.
    * Provide comprehensive API design and documentation.
* **Portability**
    * Ease of switching source of input to a different form

## Implementation Design

This API endpoint is designed for efficiency, scalability, and maintainability, with a focus on fulfilling key non-functional requirements.

**Addressing Non-Functional Requirements:**

* **Availability:**
    * **Infrastructure Redundancy:**  To ensure high availability, the application can be deployed across multiple servers with load balancing. This distributes traffic and ensures that the service remains accessible even if one server fails.
* **Scalability:**
    * **Handling Increasing Traffic:** The stateless design allows for horizontal scaling by adding more instances of the API endpoint behind a load balancer. This distributes the traffic load and ensures responsiveness even with high request volumes.
    * **Accommodating Growing Data Volumes:** The Map-Reduce paradigm enables processing large CSV files by dividing them into smaller chunks and processing them in parallel. This can be further enhanced by leveraging distributed processing frameworks like Apache Hadoop or Apache Spark. **Tradeoff:** Implementing distributed processing introduces complexity in managing chunk splitting and merging.

* **Performance:**
    * **Minimizing Response Times:** The synchronous processing approach minimizes context switching overhead, leading to faster response times. Additionally, using `bufio.Reader` for line-by-line processing reduces memory consumption and improves efficiency, especially for large files. **Tradeoff:** Using `bufio.Reader` requires handling newline characters (`\n`) explicitly, which adds to code complexity compared to using `csv.Reader`.
    * **Maximizing Throughput:**  Efficient resource utilization and the ability to scale horizontally contribute to high throughput. Asynchronous operations, such as background processing of data chunks, can be introduced to further enhance throughput.
    * **Maintaining Low Latency:**  Stateless design and optimized data processing within the same application contribute to low latency. **Tradeoff:** Data Processing logic is housed in the same application hosting the HTTP Rest API. The benefit will be lost when moving to distributed processing systems

* **Maintainability:**
    * **Clear and Organized Code Structure:** The codebase follows a modular design with clear separation of concerns. Data processing logic is decoupled from the HTTP server, promoting code reusability and maintainability.
    * **Comprehensive API Design and Documentation:**  The API design adheres to RESTful principles. Clear and concise documentation, including API specifications (e.g., OpenAPI) and code comments, is provided to facilitate understanding and maintenance.

* **Portability:**
    * **Ease of Switching Input Source:** By abstracting the data processing logic and using a generic `bufio.Reader`, the application can easily adapt to different input sources. For example, instead of reading from `request.Body`, the `bufio.Reader` could be initialized with a file or a network stream. This flexibility allows for seamless integration with various data sources.


**Key Design Choices and Rationale:**

* **Transient Data:** The implementation primarily relies on transient data, with state limited to a read-only schema constructed from the initial file read. This approach facilitates porting to a distributed system.

* **Abstracted Data Processing:** Data processing logic is decoupled from the HTTP server, accepting a generic `bufio.Reader`. This allows for potential migration to a separate microservice. Using `bufio.Reader` enables line-by-line processing of the request body, reducing memory footprint for large files.

* **Synchronous Processing:** A `TransientHierarchy` structure encapsulates the CSV data. While it supports both synchronous and concurrent processing, benchmarking revealed that synchronous processing is more efficient for this CPU-bound task.

    **Benchmarking Results:**
    ```bash
    goos: darwin
    goarch: arm64
    pkg: github.com/bochap-learning/r-project/tree
    cpu: Apple M1 Pro
    BenchmarkConcurrentExtractLargeInput-10              676    1753469 ns/op
    BenchmarkSynchronousExtractLargeInput-10            1465    806253 ns/op
    BenchmarkConcurrentExtractSmallInput-10            91137    12562 ns/op
    BenchmarkSynchronousExtractSmallInput-10          121779    10508 ns/op
    ```

    **Rationale:** The synchronous calls take lower `ns/op` regardless of file size. This is likely because the process is CPU-bound, and the overhead of context switching in coroutines outweighs the potential benefits of concurrency.

* **Pure Functions:**  Functions like `generateSchemaColumns`, `extractRecord`, and the `TreeNode` constructor are implemented as pure functions, enhancing portability across different platforms, including cloud functions.

**Map-Reduce Paradigm:**

The design follows a Map-Reduce pattern:

1. **Map:** Each CSV record is processed and mapped into a `TransientHierarchy` record (a slice of strings).

2. **Reduce:**  The `TransientHierarchy` records are reduced into a single `TreeNode` representing the final JSON output.

This paradigm allows for future scalability by splitting large files into chunks for parallel processing using a distributed Map-Reduce framework.

**File structure:**

    .
    ├── LICENSE
    ├── README.md
    ├── code_task.md
    ├── design-task-illustration.png
    ├── go.mod
    ├── io_test.go
    ├── main.go                             # http server related code
    ├── testdata
    │   ├── large_input.csv
    │   ├── large_output.json
    │   ├── small_input.csv
    │   └── small_output.json
    └── tree
        ├── benchmark_test.go               # benchmark tests for processing the transient_tree synchronously and concurrently
        ├── constants.go                    # constant and vars used in the tree package
        ├── node.go
        ├── node_test.go                    # structs, funcs and receiver funcs for slices into TreeNodes
        ├── transient_hierarchy.go          # structs, funcs and receiver funcs for processing csv records into a slices 
        └── transient_hierarchy_test.go