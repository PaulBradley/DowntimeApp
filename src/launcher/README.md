# Launcher

The launcher application is Windows a executable that health institutions can deploy to their Windows fleet.
Ideally, IT departments would deploy the executable to each users desktop folder.

When a user launches the application it gathers a users context, `Windows Username` & computer `Hostname` together
with other telemetry data like a unique `Transaction ID` before packaging the context into a signed JSON Web Token (JWT)

After creating the signed JWT, the application then launches the DowntimeApp web site in the users default web browser.
The DowntimeApp then verifies the JWT before letting the user proceed and interact with the application.

It is intended that each health institution will have their own unique launcher executable.