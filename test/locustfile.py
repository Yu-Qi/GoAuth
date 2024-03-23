from locust import HttpUser, task, constant_throughput, events
import random

class SingleRequestUser(HttpUser):
    host = "http://localhost:9030"
    wait_time = constant_throughput(5) # 5 request per second

    # @task
    # def get_product_recommendation(self):
    #     headers = {
    #         'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2FsYW5jaGVuLmNvbSIsImV4cCI6MTcxMTE3ODk0NiwianRpIjoiNTNkOTBlYTktYmI2Ni00YjkwLWJkZjEtMjZkZjgyMmY3M2I3IiwiaWF0IjoxNzExMDkyNTQ2LCJpc3MiOiJBbGFuIGNoZW4iLCJuYmYiOjE3MTEwOTI1NDYsInN1YiI6IjI3MDA2YWE2LWI5NTctNDY0My05MTI5LTQ2NWNiNWYyNDZjYyJ9.QWIQQ3VcZkgJFLrskGDEJk4tAjKSE8RaxmkszBWnEdE'
    #     }
    #     self.client.get("/products/recommendation", headers=headers)

    # each user only has 1 task
    def on_start(self):
        headers = {
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2FsYW5jaGVuLmNvbSIsImV4cCI6MTcxMTE3ODk0NiwianRpIjoiNTNkOTBlYTktYmI2Ni00YjkwLWJkZjEtMjZkZjgyMmY3M2I3IiwiaWF0IjoxNzExMDkyNTQ2LCJpc3MiOiJBbGFuIGNoZW4iLCJuYmYiOjE3MTEwOTI1NDYsInN1YiI6IjI3MDA2YWE2LWI5NTctNDY0My05MTI5LTQ2NWNiNWYyNDZjYyJ9.QWIQQ3VcZkgJFLrskGDEJk4tAjKSE8RaxmkszBWnEdE'
        }
        self.client.get("/products/recommendation", headers=headers)
    @task
    def wait_task(self):
        pass

    @events.test_start.add_listener
    def on_test_start(environment, **kwargs):
        print("Resetting statistics")
        environment.stats.reset_all()