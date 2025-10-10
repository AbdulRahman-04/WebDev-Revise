// Try and catch block: it is used for error handling in js. it is an es6 method.

// e.g: 


// try{

//     let number = 12;
//     if(number > number1){
//         console.log(`${number} is greatest`);
        
//     } else {
//         console.log(`${number1} is greatest!`);
        
//     }


// } catch(error){
 
//     console.log(error.name);
//     console.log(error.message);

// }


// try {
  
//     let myNumber = 78069;
//     if(myNumber1 === 7809){
//         console.log('yOUR NUMBER IS CORRECT');
        
//     } 
 

// } catch (error) {
    
//     console.log(error.name);
//     console.log(error.message);
    
    
// }




// object and array destructuring: this method is used for stroing array elements and object keys inside 
//                                 a variable.


// Object destructure: 
// object destructuring:

let myObj = {
    fname : 'syed',
    lname: 'samid',
    cohort: 'c25'
}

let {fname, lname, cohort} = myObj;
console.log(fname, lname, cohort);


// Array destructuring: 
let Arr = ['my', 'dream', 'job', 'is', 'to', 'become', 'software engineer!']

const [one, two, three, four, five, six, seven] = Arr
console.log(one, two, three, four, five, six, seven);

// for of and for in loop: 

// for of loop is used for only strings and arrays to read their index elements

let myStr = "AAJAO NA BOSS!";
for(let i of myStr){
    console.log(i);
    
}

let myArr = [29, true, 'suhail', 'neha', false];
for(let i of myArr){
    console.log(i);
    
}


// for in loop is used for objects only

let myObj1 = {
    fname: 'META',
    location: 'HYDERABAD'
}
for(let i in myObj){
    console.log(i, myObj[i]);
    
}

// settimeout and setinterval:

// settimeout: it is used to execute a function after a specific time and for once only!

function mySum(a, b){

     setTimeout(()=> {
        console.log(a +b);
        
     }, 4000)


}
mySum(5, 5)


// setinterval: it is used for executing the function continouslu after specific period of time.

function myName(name){

    let stopAt = setInterval(()=> {
        console.log(name);
        
    }, 1000)

    setTimeout(()=> {
        clearInterval(stopAt)
    }, 15000)

}
myName("Abdullah is name")