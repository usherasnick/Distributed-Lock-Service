# -*- coding: utf-8 -*-
from redis.sentinel import Sentinel


if __name__ == "__main__":
    sentinel = Sentinel(
        [
            ("localhost", 26379),
            ("localhost", 26380),
            ("localhost", 26381),
        ],
        sentinel_kwargs={"password": "Pwd123!@"}, 
        socket_timeout=0.1,
    )
    master = sentinel.master_for("mymaster", socket_timeout=0.1, password="sOmE_sEcUrE_pAsS")
    master.set("msg:hello", "Hello Redis!!!")
    print(master.get("msg:hello"))
