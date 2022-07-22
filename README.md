## Steps to download and make things ready
1. open a terminal
2. git clone https://github.com/HiteshRepo/ports_processor.git
3. cd ports_processor
4. go mod download
5. go mod tidy

## Steps to install task
1. Please follow: https://taskfile.dev/installation/

## Steps to run tests
1. To run all tests (unit tests and integration test): `task test`
2. To run app tests: `task ports_processor:test`
3. To run app integration tests: `task ports_processor:test:integration`
4. To run pkg tests: `task pkg:test`

## Steps to run the app and verify results
1. Command to run app:`task run:local`
2. Verify:
   1. connect to postgres instance with these details:
      1. Username: postgres
      2. Password: mysecretpassword
      3. Host: 127.0.0.1/localhost
      4. Port: 8888
      5. DBName: ports_db
   2. run the query: `select * from ports;`
