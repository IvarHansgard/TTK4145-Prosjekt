TASK 3:
What does GOMAXPROCS do? What happens if you set it to 1?
You can only run one thread at a time so the routines will run after each other 

Results from code:
We use two threads to run the increase and decrease function for the global variable i at the same time causing a race condition

RUN | RESULT 
1   | -88891
2   | 50016
3   | 66883
4   | 82282
5   | -186276

TASK 4:

DONE

TASK 5:
