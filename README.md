# odoo-golang-restapi
Alpha stages of Odoo 14 database REST API in Golang

# Purpose?
API call directly to Odoo 14 Postgres database

# TODO
- Create Middleware (Done)
- Create Auth (Partially done, currently only implement login. Next are implementing logout and updating users. For registering and delete users, impossible to implement register since odoo has a unique register user function and delete user function)
- Implement more tables or modules such as:
  - HR
  - Helpdesk
  - Expense
  - Approvals (Odoo Enterprise Only)
