db.createUser(
    {
        user: "sigma-intern",
        pwd: "sigma",
        roles: [{
            role: "readWrite",
            db: "persons"
        }   
        ]

    }
)