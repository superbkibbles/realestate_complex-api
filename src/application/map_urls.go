package application

func mapUrls() {
	router.GET("/api/complex", handler.Get)                     // Get all Compexes
	router.POST("/api/complex", handler.Create)                 // Create Complex
	router.GET("/api/complex/:complex_id", handler.GetByID)     // Get Agency By ID
	router.POST("/api/complex/:complex_id", handler.UploadIcon) // Upload the Icon
	router.PATCH("/api/complex/:complex_id", handler.Update)    // Update Agency
	router.GET("/api/complex/search/s", handler.Search)         // Search for properties
	router.DELETE("/api/complex/:complex_id", handler.DeleteIcon)
}
