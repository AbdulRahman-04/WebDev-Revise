import config from "config"
import mongoose from "mongoose"

const dbUrl = config.get("DB_URL")

async function DbConnect() {
    try {
 
        await mongoose.connect(dbUrl)
        console.log("MONGODB CONNECTED SUCCESSFULLY✅");
        
         

        
    } catch (error) {
        console.log("Couldn't Connect to mongoDB");
        
    }
}

DbConnect()