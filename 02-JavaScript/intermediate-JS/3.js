/*
  Functions : A function in programming is a block of code that does a specific task, 
             and you can reuse it whenever you need.
*/

// simple function
// function with parameter and arguements
// function with return
// default parameter function
// function expression
// anonymous function
// IIFE
// Arrow function

// simple function: 

function SayHello(){
  console.log("Hello from simple function");
}

SayHello()

// function with parameter and arguements
function MyParam(var1, var2){
  console.log(var1 + var2);
}

MyParam(10, 10)

// function with return
function MyReturn(age){
  
  if (age >= 18){
    console.log("You can drive");
    return
  } 

  console.log("You cant drive");
  

}

MyReturn(7)


// default parameter function

function DefaultParam(n1 = 12, n2= 15){
  console.log(n1, n2);
}

DefaultParam()


// Hoisting (calling a function even before its declared is possible in above functions)


// function expression: Dekh bro, function expression basically ek function ko variable ke andar store karna hota hai
// Hoisting nahi hoti (yaani upar call nahi kar sakte pehle).
// Anonymous function (no name) ya named function dono ho sakte.

const variable = function(value){
  return value
}

let val = variable("Hello bro")
console.log(val);


// Anonymous functions: these functions are basically functions which have no name and the example of
//  these function are function expression and arrow function.


// Arrow function : these are an important functions used in advance js and node js and its syntax is 
//                   similar to function expression.

const MyRollNo = (rollNo) => {
  return rollNo
}

const value = MyRollNo(5049);
console.log(value);



// Hoisting: a function can be called anywhere in the code even before its written this is called as
//           hoisting. hoisting can't be done for anonymous functions(function expression and arrow functions).



SayHey("hey")

function SayHey(variable){
  console.log(variable);
  
}


// Recursion: a function called inside itself is called recursion. because of recursion the function
//           gets into infinite loop and  its output gets executed infinte times until stopped manually

function loop(value){
  console.log(value);

  loop()
  
}

loop(10)