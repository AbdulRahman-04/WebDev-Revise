// function setTimer(a, b){
 
//     setTimeout((x)=> {
//         console.log(a + b);
        
//     }, 4000)

// }
// setTimer(2, 5)

// function myName(name){

//     let stopAt = setInterval(()=> {
//         console.log(name);
        
//     }, 1000)

//     setTimeout(()=> {
//         clearInterval(stopAt)
//     }, 15000)

// }
// myName("Abdul Rahman")

function SayName(name){
    let stopAt = setInterval(()=>{
        console.log(name);
        
    }, 1000)

    setTimeout(()=>{
        clearInterval(stopAt)
    }, 5000)
}

SayName("Syed Bhai")