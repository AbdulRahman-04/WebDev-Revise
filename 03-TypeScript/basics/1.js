// Basics Types in typescript: 
// 1. Primitive types : number string boolean null undefined bigint symbol
// 2. Arrays
// 3. Tuples
// 4. enum 
// any, unknown, void, null, undefined, never
/* Primitive types in typescript:  there are 3 primitive types in ts, they are: number, string and boolean.

   just like in js these primitive variables are declared.

   in typescript, u can do type annotation, which means telling compiler at time of
   decalring variable its data type

 */
// primitive type: 
var userName = "rahman";
var age = 21;
var isAlive = true;
var isEmpty = null;
var notDefined = undefined;
console.log(userName, age, isAlive, isEmpty, notDefined);
// /* Arrays in TypeScript: unlike js array, ts array has power to declare what
//  type of value is going to be stored inside the array. e.g: number or string
// normal array 
var myArr = ["hi", "hey"];
console.log(myArr);
// ts array : 
var array = [true, false, true, false];
console.log(array);
var numbers = [49, 69, 89];
console.log(numbers);
var strings = ["hi", "kysa", "h?"];
console.log(strings);
/*
  Tuples : tuple is a special type of array,
   which stores fixed size and specific data values in itself.

  let arr = [string, boolean, number] = ["hey", true, "go"]

*/
var tupleArr = ["hi", true, 0];
console.log(tupleArr);
var tupleArr2 = ["h", "i", true, false, 0];
console.log(tupleArr2);
// ENUM: An enum is a way to store a fixed set of values under a name. 
//       It makes the code clean and easy to read.
var UserRoles;
(function (UserRoles) {
    UserRoles["ADMIN"] = "Admin";
    UserRoles["USER"] = "user";
    UserRoles["SUPERADMIN"] = "superAdmin";
})(UserRoles || (UserRoles = {}));
var StatusCodes;
(function (StatusCodes) {
    StatusCodes[StatusCodes["ERROR"] = 404] = "ERROR";
    StatusCodes[StatusCodes["NOINTERNET"] = 500] = "NOINTERNET";
})(StatusCodes || (StatusCodes = {}));
console.log(UserRoles.ADMIN, UserRoles.USER, UserRoles.SUPERADMIN);
console.log(StatusCodes.ERROR, StatusCodes.NOINTERNET);
// ANY , UNKNOWN , VOID, UNDEFINED, NULL ETC:
// ANY : If a variable is declared without specifying its type and without assigning a value, TypeScript treats it as any. 
//       This means the variable can hold any type of value, which removes TypeScript’s type safety and is generally not recommended.
var myVar; // this is an any type variable in ts
var myVar3 = 2025; // this is also any type variable
console.log(myVar3);
var a;
// a = true   here i can't give value to a of any other dt except number bcos i already defined what dt value it'll hold in future in line 110
var b; // here the value and data type of b can be of ANY data type in TS. 
// Unknown: it is a special TypeScript type that can store any value, but unlike any,
//         it does not allow operations (like arithmetic or method calls) without proper type checking..
var d = "hello";
d = d.toUpperCase();
console.log(d);
var x = 9;
if (typeof x === "string") {
    console.log(x);
}
// any datatype mtlab ek aisa variable jiska type and value kuch b ho skta, which is a threat to our typeSafety code, so its always better
// to use unknown as type to a variable which we aren't sure what type of value is gonna be stored inside it and what will be its type.
/* Void: it is a special type which tells that what datatype value
        is getting returned inside a function. if no value is getting retrned
        
        void ka matlab hai “koi value return nahi ho rahi”.

        Matlab function sirf execute hoga, lekin kuch return nahi karega.

        */
function sayName() {
    console.log("rahman bhai");
}
sayName();
function sayNumber(myVaris) {
    return myVaris;
}
var value = sayNumber(5049);
console.log(value);
// undefined: it literally means a variable is declared but it has been not assigned with any value.
var myValueis;
console.log(myValueis);
// never: use karte ho jab function kabhi normal tarike se end nahi hoga. 
//        Ya toh error throw karega ya infinite loop me fasa rahega.
function throwError(message) {
    throw new Error(message);
}
// throwError("error")
function infiniteLoop() {
    while (true) {
        console.log("Running forever...");
    }
}
// infiniteLoop()
// Type inference and type annotations: 
// type inference: it means when we declare a varaible without defing its type. ts automatically checks the data type of variabke its called 
//                 the type inference
var myCollege = "DCET";
console.log(myCollege, typeof myCollege);
//  type annotations: it means when we declare a variable while also defining its data type.
var getVar = true;
console.log(getVar, typeof getVar);
