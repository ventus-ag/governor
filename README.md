# Governor

Governor is a service that will automatically manage your resources, monitor and offer you the most profitable use of your resources in accordance with your business plan.

## How it works

Every thirty minutes Governor getting three type of resources (CPU Cores, RAM, Instances) of all projects. Then, the service checks how many resources client have on project and how many resources client use. And if currently used resources of project have passed treshold, that equal sixty percent, Governor will send message on project's owner e-mail.

## Get Started using Governor

To get started using Governor you need:

- Create project in [Ventus Cloud portal](https://portal.ventuscloud.eu/) 
  Check that you have this project in openstack [Openstack Ventus portal](https://cloud.vstack.ga/identity/)
  Project name must include *Name - id*, for example: `Test Ventus - 123456`

- Set up your project to use one of the resources (CPU Cores, RAM or Instances) to 60%. 

- Within half an hour you will receive a letter with a notification that you have exceeded the threshold of 60% in the use of a particular resource.
  If your project pass the threshold in two or three type of recourses. You will get separate letter for each resource.

  In the case that your project passes the threshold for using less than 60% of recourse, for example, you will not increase the number of Instances, but increase their size, the information that you received a notification will be updated. And when you pass the threshold of 60% again - you will receive a notification letter.
