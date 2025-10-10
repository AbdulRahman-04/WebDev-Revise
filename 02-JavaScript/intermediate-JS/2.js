// Operators in js :

// there are several types of operators in js :

// 1. Arithmetic operators : +, -, *, /, %, ++, --

var a = 10;
var b = 20;

console.log(a + b, a - b, a * b, a % b, ++a, --b);

// 2. Assignment operators : =, +=, -=, *=, /=, %=

var c = 35;

c += 5;
console.log(c);

var d = 50;

d %= 30;
console.log(d);

// Comparison operators : ==, ===, !=, !==, >, <, >=, <=

let ab = "123";
let ba = 123;

//    console.log(ab == ba);
//    console.log(ab === ba);
//    console.log(ab != ba);
//    console.log(ab !== ba);

// Logical operators : &&, ||, !
let x = 10;
let y = 20;

//   console.log(x <= y && y >= x);
console.log(x >= y || y > x);
console.log(!(x > y));

// Ternary Operator: used for single line conditions.

let age = 17;

age >= 18
  ? console.log("You can drive")
  : console.log("You are a minor, u cant drive");

// Conditionals in Javascript: the conditional statements in javascript are of 4 types:

/*
 --> If statement
 --> IF-ELSE statement
 --> IF ELSE-IF ELSE statement
 --> NESTED IF AND IF ELSE statement
*/

// IF  statement: it is used for single line statements only:

let number = 12;

// if (number <= 15) {
//    console.log("number is less than or equal to 15");

// }

// IF-ELSE  statement:

if (number <= 15) {
  console.log("number is less than or equal to 15");
} else {
  console.log("number is greater than 15");
}

// // IF..ELSE-IF..IF statement:
let myCohort = "C24";

if (myCohort === "A24") {
  console.log("This is not my cohort");
} else if (myCohort === "B24") {
  console.log("This is not my cohort");
} else if (myCohort === "C24") {
  console.log("This is my cohort bro😭🤌");
} else {
  console.log("None of the above");
}



// // nested if else :

let temp = 6;

if (temp){
  
    console.log("Checking temperature");

    if (temp > 0) {
        console.log("Temperature is greater than 0 degree");
        
    } else if (temp >= 10) {
        console.log("Temperature is greater than or equal to 10 degree");
        
    }  else if (temp >= 20) {
        console.log("Temperature is greater than or equal to 20 degree");
        
    } else if (temp >= 30) {
        console.log("Temperature is greater than or equal to 30 degree");
        
    } 
} else {
    console.log("Temperature is less than 0 degree");
    
}


