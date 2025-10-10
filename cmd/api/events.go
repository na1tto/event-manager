package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// this is the handler for the events operations
// handlers (or controllers in other contexts), are responsible for receiving the request and responding

// createEvent - Creates a event using the bearer token
//
// @Summary 				Creates a event
// @Description 		Creates a event in the database using a bearer token
// @Tags 						Events
// @Accept 					json
// @Produce 				json
// @Param 					credentials body				database.Event	true	"Event informations"
// @Success 				200 					{object} 	database.Event
// @APIResponses(code = 400, message = "Bearer token is required", response = err.Error())
// @ApiResponse(code = 400, message = "Bad request, adjust before retrying", response = err.Error())
// @ApiResponse(code = 500, message = "Internal server error, failed to create event", response = err.Error())
// @Security				BearerAuth
// @Router /api/v1/events [POST]
func (app *application) createEvent(c *gin.Context) {
	var event database.Event //empty Event struct

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.getUserFromContext(c)
	event.OwnerId = user.Id

	err := app.models.Events.Insert(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// getEvent return all events
//
// @Summary Returns a single event
// @Description Returns a single event
// @Tags Events
// @Accept json
// @Produce json
// @Param id	path	int	true "event ID"
// @Success 200 {object} object{error=string}
// @Router /api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not Found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// getEvents return all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// updateEvent - updates a event using the bearer token
//
// @Summary 				Updates a event
// @Description 		Updates a event in the database using a bearer token
// @Tags 						Events
// @Accept 					json
// @Produce 				json
// @Param						id path int true "event ID"
// @Param 					credentials body				database.Event	true	"Updated event informations"
// @Success 				200 					{object} 	database.Event
// @Failure 400 {string} string = "Error: Invalid event ID"
// @Failure 404 {string} string = "Error: Event not found"
// @Failure 403 {string} string = "Error: You are not authorized to update this event"
// @Failure 500 {string} string = "Error: Failed to retrieve event"
// @Security				BearerAuth
// @Router /api/v1/events/{id} [PUT]
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
	}

	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent - deletes a event using the bearer token
//
// @Summary 				Deletes a event from the database
// @Description 		Deletes a event from the database using a bearer token
// @Tags 						Events
// @Accept 					json
// @Produce 				json
// @Param						id path int true "event ID"
// @Success 				204
// @Failure 400 {string} string = "Error: Invalid event ID"
// @Failure 404 {string} string = "Error: Event not found"
// @Failure 403 {string} string = "Error: You are not authorized to delete this event"
// @Failure 500 {string} string = "Error: Failed to retrieve event"
// @Security				BearerAuth
// @Router /api/v1/events/{id} [DELETE]
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// addAttendeeToEvent - adds a attendee to a event
//
// @Summary 				Adds a attendee to a event in the databse
// @Description 		Adds a attendee to a event in the database (needs authetication)
// @Tags 						Attendees
// @Accept 					json
// @Produce 				json
// @Param						id path int true "event ID"
// @Param						userId path int true "attendee to be added"
// @Success 				201	 {object} database.Attendee
// @Failure 400 {string} string = "Error: Invalid event ID"
// @Failure 403 {string} string = "Error: You are not authorized to delete this event"
// @Failure 404 {Error}  Error  = err.Error()
// @Failure 409 {string} string = "Error: Attendee already exists"
// @Failure 500 {string} string = "Error: Failed to retrieve event"
// @Security				BearerAuth
// @Router /api/v1/events/{id}/attendees/{userId} [POST]
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add ateendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// getAllAttendeesForEvent - gets the attendees list from a event
//
// @Summary 				Gets the attendees list from a event
// @Description 		Gets the attendees list from a event
// @Tags 						Attendees
// @Accept 					json
// @Produce 				json
// @Param						id path int true "event ID"
// @Success 				200 {object} []database.User
// @Failure 400 {Error} Error = err.Error()
// @Failure 500 {Error} Error = err.Error()
// @Router /api/v1/events/{id}/attendees/ [GET]
func (app *application) getAllAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetAttendeesByEvent - deletes a attendee from a event
//
// @Summary 				Deletes a attendee from a event
// @Description 		Deletes a attendee from a envent through the event and attendee's ID
// @Tags 						Attendees
// @Accept 					json
// @Produce 				json
// @Param						id path int true "event ID"
// @Param						userId path int true "user ID"
// @Success 				204
// @Failure 400 {Error} Error = err.Error()
// @Failure 401 {Error} Error=err.Error()
// @Failure 403 {Error} Error = err.Error()
// @Failure 404 {Error} Error = err.Erro()
// @Failure 500 {Error} Error = err.Error()
// Security 				BearerAuth
// @Router /api/v1/events/{id}/attendees/{userId} [DELETE]
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete an attendee from event"})
	}

	err = app.models.Attendees.Delete(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// getEventsByAttendee - gets the eventes per attendee
//
// @Summary 				Gets the list of events a user is attending
// @Description 		Gets tge list of events a user is attending using it's ID
// @Tags 						Attendees
// @Accept 					json
// @Produce 				json
// @Param						id path int true "user ID"
// @Success 				200
// @Failure 400 {Error} Error = err.Error()
// @Failure 500 {Error} Error = err.Error()
// Security 				BearerAuth
// @Router /api/v1/attendees/{id}/events [GET]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
