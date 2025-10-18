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
let userName:string = "rahman"    
let age : number = 21
let isAlive: boolean = true
let isEmpty:null = null
let notDefined: undefined = undefined


console.log(userName, age, isAlive, isEmpty, notDefined);


   

// /* Arrays in TypeScript: unlike js array, ts array has power to declare what
//  type of value is going to be stored inside the array. e.g: number or string

// normal array 
let myArr = ["hi", "hey"]
console.log(myArr);

// ts array : 
let array:Boolean[] = [true, false , true, false]
console.log(array);

let numbers: number[] = [49, 69, 89]
console.log(numbers);


let strings : string[] = ["hi", "kysa", "h?"]
console.log(strings);


/*
  Tuples : tuple is a special type of array,
   which stores fixed size and specific data values in itself.

  let arr = [string, boolean, number] = ["hey", true, "go"] 

*/ 

let tupleArr: [string, boolean, number] = ["hi", true, 0]
console.log(tupleArr);

let tupleArr2 : [string, string, boolean, boolean, number] = ["h", "i", true, false, 0]
console.log(tupleArr2);


// ENUM: An enum is a way to store a fixed set of values under a name. 
//       It makes the code clean and easy to read.

enum UserRoles {
    ADMIN = "Admin",
    USER = "user",
    SUPERADMIN = "superAdmin"
}

enum StatusCodes{
    ERROR = 404,
    NOINTERNET = 500
}

console.log(UserRoles.ADMIN, UserRoles.USER, UserRoles.SUPERADMIN);
console.log(StatusCodes.ERROR, StatusCodes.NOINTERNET);


// ANY , UNKNOWN , VOID, UNDEFINED, NULL ETC:

// ANY : If a variable is declared without specifying its type and without assigning a value, TypeScript treats it as any. 
//       This means the variable can hold any type of value, which removes TypeScript’s type safety and is generally not recommended.

let myVar; // this is an any type variable in ts
let myVar3:any = 2025 // this is also any type variable
console.log(myVar3);


let a: number;

// a = true   here i can't give value to a of any other dt except number bcos i already defined what dt value it'll hold in future in line 110


let b; // here the value and data type of b can be of ANY data type in TS. 


// Unknown: it is a special TypeScript type that can store any value, but unlike any,
//         it does not allow operations (like arithmetic or method calls) without proper type checking..

let d : any = "hello"

d = d.toUpperCase()
console.log(d);


let x:unknown = 9

if (typeof x === "string"){
    console.log(x);
    
}


// any datatype mtlab ek aisa variable jiska type and value kuch b ho skta, which is a threat to our typeSafety code, so its always better
// to use unknown as type to a variable which we aren't sure what type of value is gonna be stored inside it and what will be its type.


/* Void: it is a special type which tells that what datatype value 
        is getting returned inside a function. if no value is getting retrned 
        
        void ka matlab hai “koi value return nahi ho rahi”.

        Matlab function sirf execute hoga, lekin kuch return nahi karega.

        */

function sayName():void {
    console.log("rahman bhai");
}        

sayName()

function sayNumber(myVaris:number):number{
    return myVaris;
}

let value:number = sayNumber(5049)
console.log(value);


// undefined: it literally means a variable is declared but it has been not assigned with any value.

let myValueis:undefined;
console.log(myValueis);



// never: use karte ho jab function kabhi normal tarike se end nahi hoga. 
//        Ya toh error throw karega ya infinite loop me fasa rahega.

function throwError(message: string): never {
    throw new Error(message);
}

// throwError("error")

function infiniteLoop(): never {
    while(true) {
        console.log("Running forever...");
    }
}
// infiniteLoop()



// Type inference and type annotations: 

// type inference: it means when we declare a varaible without defing its type. ts automatically checks the data type of variabke its called 
//                 the type inference

let myCollege = "DCET"
console.log(myCollege, typeof myCollege);

//  type annotations: it means when we declare a variable while also defining its data type.
let getVar : boolean = true;
console.log(getVar, typeof getVar);


//  interfaces and type alises: 

// interface: an interface is like a rulebox used for objects , as it tells what properties an object must have and also their types.

interface myDetails{
    name: string
    age:number
    rollno: boolean
}

function userDetails(details: myDetails){
   
    console.log(details);
    


}

// userDetails({name:"abdul", age:21, rollno:true})


// extending interfaces : 

interface Admin extends myDetails {
    admin: boolean
}

function isAdmin(obj : Admin){
    console.log(obj);
    
}

isAdmin({name:"abdul", age:21, rollno:true, admin:true})


// type aliases : creating own custom datatype

type lod = string | null | boolean 

let myName: lod = "fahad"
console.log(myName);


type myCustomType = string | number | boolean | null 

let ab : myCustomType = 29
console.log(ab);


// union type : | or operator from js 

let bc: string | number | boolean = true
console.log(bc);


// intersection type 
type myType = lod & {
    password: string,
    age: number
}
