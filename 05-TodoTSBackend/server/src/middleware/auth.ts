import jwt from "jsonwebtoken";
import config from "config";
import { Request, Response, NextFunction } from "express";

const KEY: string = config.get<string>("JWT_KEY");

interface AuthRequest extends Request {
  user?: any;
}

const authMiddleware = (req: AuthRequest, res: Response, next: NextFunction): void => {
  const authHeader = req.headers["authorization"];

  if (!authHeader) {
    res.status(401).json({ msg: "No token provided ❌" });
    return;
  }

  const token = authHeader.split(" ")[1];

  try {
    const decoded = jwt.verify(token, KEY);
    req.user = decoded;
    next();
  } catch (error) {
    console.error("JWT verification failed:", error);
    res.status(403).json({ msg: "Invalid or expired token ❌" });
  }
};

export default authMiddleware;