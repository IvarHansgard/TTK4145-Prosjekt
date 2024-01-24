// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>
pthread_mutex_t lock;
int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    pthread_mutex_lock (&lock);
    // TODO: increment i 1_000_000 times
    for(int j=0; j<1000000; j++){
        i++;
    }    
    pthread_mutex_unlock(&lock);
    return NULL;
}

void* decrementingThreadFunction(){
    pthread_mutex_lock (&lock);
    // TODO: decrement i 1_000_000 times
    for (int k=0; k<1000000;k++){
        i--;
    }
    pthread_mutex_unlock(&lock);
    return NULL;
}


int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_t thread1, thread2;

    pthread_create(&thread1, NULL, incrementingThreadFunction, NULL);
    pthread_create(&thread2, NULL, decrementingThreadFunction, NULL);

    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`  
    pthread_join(thread1, NULL);
    pthread_join(thread2, NULL);  
    pthread_mutex_destroy(&lock);
    printf("The magic number is: %d\n", i);
    return 0;
}
