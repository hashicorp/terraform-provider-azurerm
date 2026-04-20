# Example: Logic App Standard with Content Share Disabled

This example provisions a Logic App Standard with `content_share_force_disabled = true`, which suppresses the `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` and `WEBSITE_CONTENTSHARE` app settings.

This is useful when deploying to an App Service Environment (ASE) where the recommendation is to use the ASE internal storage account instead of an external Azure Files share.
