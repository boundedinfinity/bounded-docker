package docker

func CompareId(id1, id2 string) bool {
	if len(id1) == 0 || len(id2) == 0 {
		return false
	}

	if len(id1) < len(id2) {
		return id2[:len(id1)] == id1
	}

	if len(id1) > len(id2) {
		return id1[:len(id2)] == id2
	}

	return id1 == id2
}
