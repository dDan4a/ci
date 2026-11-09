RANK = {"critical": 0, "warning": 1, "info": 2}


def compare_severity(a: str, b: str) -> int:
    return RANK.get(a, RANK["info"]) - RANK.get(b, RANK["info"])
