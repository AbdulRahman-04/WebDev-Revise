import mongoose from "mongoose";
import config from "config"


const db_url:string = config.get<string>("DB_URL")

async function dbConnect():Promise<void> {

    try {
      
        
        await mongoose.connect(db_url)
        console.log("MONGO DB CONNECTED SUCCESSFULLY!✅");
        


        
    } catch (error) {
        console.log(error);
        
    }
     
}

dbConnect()