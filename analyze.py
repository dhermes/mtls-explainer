import binascii


VERBS = ("Read", "Write")


def main():
    with open("info.txt", "r") as file_obj:
        lines = file_obj.read()

    for row in lines.strip().split("\n"):
        verb, size, content = row.split("|")
        if verb not in VERBS:
            raise ValueError("Invalid row", row)
        if len(content) != 2 * int(size):
            raise ValueError("Invalid row", row)

        raw_content = binascii.unhexlify(content.encode("ascii"))
        print(verb)
        print(raw_content)


if __name__ == "__main__":
    main()
