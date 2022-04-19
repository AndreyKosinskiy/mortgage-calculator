# mortgage-calculator
Mortgage Calculator is web application where users can create banks and calculate mortgage payments using one of these bank’s settings.

## Getting Started
This is an example of how you may give instructions on setting up your project locally. To get a local copy up and running follow these simple example steps.

1. Clone the repo:
    ```sh
       git clone https://github.com/AndreyKosinskiy/mortgage-calculator.git
    ```
2. Run it!

      2.1 Fill .env file or delete .example extention in .env.example file.
      
      2.2 On local machine:
      ```sh
         go mod download
         make migrate # require postgresql
         make build_
         make run_
      ```
      2.3 In Docker:
      ```sh
         make run
      ```
## Roadmap

- [x] Add Docker
- [x] Use inmemory storage
- [x] Use html/template package
- [x] Use PostgreSQL
- [x] Write test for func calcMMP( month martgage payment )
- [x] Add success month martagage payment to localStorage
- [ ] Wirite API
- [ ] Change html/template to React

