# Mattermost Customer Plugin

The goal is to make customer information and environment details that exist within internal systems to become more available within Mattermost in a secure and accessible way.

## Components of the Plugin (goals)

1. When a support packet is uploaded anywhere it auto parses that packet, saves the data from it, and returns a summary of the packet along with proactive findings. 
2. A sidebar exists where you can view and query the information. 
    - This should allow updating the fields
    - Includes a history of packets that have been uploaded already
    - Able to diff between dates
3. (longterm) proactively monitor zendesk for updates containing a packet. Download to Mattermost, parse, and send the ticket a private note with proactive findings. 
