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

// loop(10)


// -----------------------------------------------


// Objects: objects in js can store values of more than one data type.

let myObj = {
  name : "Rahman",
  email: "rahman69@gmail.com",
  no: 8186978069
}

console.log(myObj.name, myObj.email, myObj.no);

myObj.name = "Abdul rahman"

console.log(myObj.name, myObj.email, myObj.no);

myObj.address = "talab katta"
console.log(myObj.name, myObj.email, myObj.no, myObj.address);

delete myObj.no
console.log(myObj.name, myObj.email, myObj.address);



// there are mainly 11 methods in array which are very useful.

// 1.object.keys(): 


let newObj = {
    myRollno: 66,
    myClass: 'ece- 4th year!'
}

console.log(Object.keys(newObj));


// 2.object.values(): 

let newObj1 = {
    myRollno: 66,
    myClass: 'ece- 4th year!'
}

console.log(Object.values(newObj1));


// 3.object.Entries(): 
let newObj3 = {
    myRollno: 24,
    myClass: ' 4th year!'
}
console.log(Object.entries(newObj3));



// 4.Object.assign(): 
let newObj4 = {
    myRollno: 66,
    myClass: 'ece- 4th year!'
}

let newObj5 = {
    myStreet: 'chandulal',
    myLocation: 'baradari'
}


console.log(Object.assign(newObj4, newObj5));


// 5.Object.Create(): 

let myFilms = {
    movie1: 'DDLJ',
    movie2: 'Raees'
}

let srkMovies = Object.create(myFilms);
console.log(srkMovies.movie1);


// 6.Object.Freeze(): 

let newObj7 = {
    myRollno: 66,
    myClass: 'ece- 4th year!'
}

console.log(Object.freeze(newObj7));

newObj7.myWork = 'hw krna';
console.log(newObj7);

// object.isfrozen(): checks if an object is freeze or not!

// 7.Object.fromEntries(): 

let arr = [['myFname', 'Syed'], ['myLname', 'Abdul Rahman'], ['myLocation','Bahadurpura']]

console.log(Object.fromEntries(arr));


// 8.Object.is(): it is basically === like method 

let a = 45;
let b = 45;
// console.log(Object.is(a ,b));

// 8.Object.Seal(): 
let newObj9 = {
    myRollno: 66,
    myClass: 'ece- 4th year!'
}

Object.seal(newObj9);
newObj9.myCollege = "DCET";
// console.log(newObj9);
newObj9.myRollno = 5049;
// console.log(newObj9);

// object.isseal(): checks if an object is sealed or not: 


// 10.Object.toString(): converts numbers into another system
let smd = 123;
let abd = smd.toString(2)
console.log(abd);