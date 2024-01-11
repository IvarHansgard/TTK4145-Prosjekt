Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> *Your answer here*
concurency is when the cpu switches between running tasks and parallelism is when task run at the same time

What is the difference between a *race condition* and a *data race*? 
> *Your answer here* 
a race condition is when multiple programs needs to access the same data and the execution of the programs can change depending on which program gets the data first

a data race is almost the same but atleas one of the programs is writing  to the data withouth any locks to controll the acces to the data, causing the results to be diffrent depending on which program gets to the data first 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> *Your answer here* 
The scheduler determines what tasks should be executed first on the cpu, it does it by using diffrent algorithms like "fifo, "edo" and "sjf" 

### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> *Your answer here*
You can work on diffrent data at the same time 

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> *Your answer here*


Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> *Your answer here*
Both, it makes it easier to do multiple tasks at the same time but it can be hard to implement in code / get your head around it.

What do you think is best - *shared variables* or *message passing*?
> *Your answer here*
message passing, because its more clean

