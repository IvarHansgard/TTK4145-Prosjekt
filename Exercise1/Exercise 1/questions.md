Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> concurrency er når CPU bytter mellom oppgavene som må løses. Parallelism er når de kjører samtidig. 

What is the difference between a *race condition* and a *data race*? 
> race condition er når flere programmer trenger samm data, og det er tilfeldig hvilket program som får dataen først og kjører. Data race er når minst ett program må skrive til dataen uten at det låser tilgangen, som forårsaker at resultatet kan bli forskjellig avhenig av hilket program som kommer til dataen først. 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> scheduler velger hvilken thread som skal kjøre, ved å bruke algoritmer som fifo, edo, sjf


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> Da kan man jobbe med forskjellig data samtidig. Programmet kan også kjøre fortere

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> 

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Begge. Fordi man kan gjøre flere oppgaver samtidig, men det kan være vanskeligere å implementere og forstå

What do you think is best - *shared variables* or *message passing*?
> message passing, fordi det er sikrere


