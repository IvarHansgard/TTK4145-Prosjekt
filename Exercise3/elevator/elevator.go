package elevator

import "fmt"
//egt enum, så kanskje bytte denne syktaksen
type ElevatorBehaviour struct{
    EB_Idle bool
    EB_DoorOpen bool
    EB_Moving bool
} 
//egt enum i c, så kanskje bytte denne syntaksen
type ClearRequestVariant struct {
    // Assume everyone waiting for the elevator gets on the elevator, even if 
    // they will be traveling in the "wrong" direction for a while
    CV_All,
    
    // Assume that only those that want to travel in the current direction 
    // enter the elevator, and keep waiting outside otherwise
    CV_InDirn,   
}
type Elevator struct {
    floor int
    dirn Dirn   
    requests[N_FLOORS][N_BUTTONS] int                
    behaviour ElevatorBehaviour   
}
type config struct {
    clearRequestVariant ClearRequestVariant 
    doorOpenDuration_s double              
}

func eb_toString(eb ElevatorBehaviour) string{
    switch eb{
    case EB_Idle:
        return "EB_Idle"
    case EB_DoorOpen:
        return "EB_DoorOpen"
    case EB_Moving:
        return "EB_Moving"
    default:
        return "EB_UNDEFINED"
    }
}

func elevator_print(es Elevator){
    fmt.Println("  +--------------------+\n")
    fmt.Printf(
        "  |floor = %-2d          |\n"
        "  |dirn  = %-12.12s|\n"
        "  |behav = %-12.12s|\n",
        es.floor,
        elevio_dirn_toString(es.dirn),
        eb_toString(es.behaviour)
    )
    fmt.Println("  +--------------------+\n")
    fmt.Println("  |  | up  | dn  | cab |\n")
    for f:=N_FLOORS-1; f>=0; f-- {
        fmt.Println("  | %d", f)
        for btn:=0 ; btn<N_BUTTONS; btn++{
            if (f==N_FLOORS-1 && btn==B_HallUp) || (f==0 && btn==B_HallDown){
                fmt.Print("|     ")
            } else{
                if es.requests[f][btn] !=0 {
                    fmt.Print("|  #  ")
                }else {
                    fmt.Print("|  -  ")
                } 
            }
        }
        fmt.Println("|\n")
    }
    fmt.Println("  +--------------------+\n")
}

// i C:
/*void elevator_print(Elevator es){
    printf("  +--elevio------------------+\n");
    printf(
        "  |floor = %-2d          |\n"
        "  |dirn  = %-12.12s|\n"
        "  |behav = %-12.12s|\n",
        es.floor,
        elevio_dirn_toString(es.dirn),
        eb_toString(es.behaviour)
    );
    printf("  +--------------------+\n");
    printf("  |  | up  | dn  | cab |\n");
    for(int f = N_FLOORS-1; f >= 0; f--){
        printf("  | %d", f);
        for(int btn = 0; btn < N_BUTTONS; btn++){
            if((f == N_FLOORS-1 && btn == B_HallUp)  || 
               (f == 0 && btn == B_HallDown) 
            ){
                printf("|     ");
            } else {
                printf(es.requests[f][btn] ? "|  #  " : "|  -  ");
            }
        }
        printf("|\n");
    }
    printf("  +--------------------+\n");
}*/

func elevator_uninitialized() Elevator{
    return Elevator{
        floor:  -1,
        dirn:  D_Stop,
        behaviour:  EB_Idle,
        config: Config{
            clearRequestVariant:  CV_All,
            doorOpenDuration_s:  3.0,
        },
    }
}

