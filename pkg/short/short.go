package shortImg

func ShortImg(img string) string{
	if img == "" {
		return ""
	}

	if len(img) > 30 {
        return img[:30] + "..."
    }
    
    return img
}