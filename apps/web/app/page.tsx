export default function HomePage() {
  return (
    <main className="container mx-auto px-4 py-8">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">
          Welcome to Aegis Trust Ecosystem
        </h1>
        <p className="text-lg text-gray-600 mb-8">
          A Web3-driven sex-positive trust ecosystem built on verified human interactions
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">Review Library</h2>
            <p className="text-gray-600">
              Verified human reviews for adult products and content
            </p>
          </div>
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">Event Platform</h2>
            <p className="text-gray-600">
              Safe, verified events with digital consent
            </p>
          </div>
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">Reputation System</h2>
            <p className="text-gray-600">
              Build trust through consensual interactions
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}